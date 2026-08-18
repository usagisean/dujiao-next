#!/usr/bin/env bash
set -Eeuo pipefail

# Root-owned deployment entrypoint invoked by the dedicated zxdeploy user.
# It only accepts a staged binary from zxdeploy's incoming directory, verifies
# its checksum and architecture, snapshots the current binary/database, then
# performs an atomic replacement with health-check rollback.

readonly APP_DIR="/opt/dujiao-next"
readonly APP_BIN="${APP_DIR}/dujiao-next"
readonly DB_FILE="${APP_DIR}/db/dujiao.db"
readonly INCOMING_DIR="/home/zxdeploy/incoming"
readonly BACKUP_ROOT="/var/backups/dujiao-next/deploys"
readonly SERVICE="dujiao-next.service"
readonly HEALTH_URL="http://127.0.0.1:18082/health"

die() {
  printf 'deploy error: %s\n' "$*" >&2
  exit 1
}

[[ $# -eq 3 ]] || die "usage: zx-shop-deploy <artifact> <sha256> <commit>"

artifact="$(readlink -f -- "$1")"
expected_sha="${2,,}"
commit="$3"

[[ "$artifact" == "${INCOMING_DIR}/dujiao-next."* ]] || die "artifact path is outside incoming directory"
[[ -f "$artifact" && ! -L "$artifact" ]] || die "artifact must be a regular file"
[[ "$(stat -c '%U' "$artifact")" == "zxdeploy" ]] || die "artifact owner must be zxdeploy"
[[ "$expected_sha" =~ ^[0-9a-f]{64}$ ]] || die "invalid sha256"
[[ "$commit" =~ ^[0-9a-f]{40}$ ]] || die "invalid commit sha"

actual_sha="$(sha256sum "$artifact" | awk '{print $1}')"
[[ "$actual_sha" == "$expected_sha" ]] || die "artifact checksum mismatch"

artifact_size="$(stat -c '%s' "$artifact")"
(( artifact_size >= 10000000 && artifact_size <= 200000000 )) || die "unexpected artifact size"
file "$artifact" | grep -Eq 'ELF 64-bit.*x86-64' || die "artifact is not a Linux x86-64 executable"

exec 9>/run/lock/zx-shop-deploy.lock
flock -n 9 || die "another deployment is already running"

timestamp="$(date -u +%Y%m%dT%H%M%SZ)"
short_commit="${commit:0:12}"
snapshot_dir="${BACKUP_ROOT}/${timestamp}-${short_commit}"
new_bin="${APP_DIR}/.dujiao-next.new.${short_commit}"

install -d -m 0700 -o root -g root "$snapshot_dir"
install -m 0755 -o root -g root "$APP_BIN" "${snapshot_dir}/dujiao-next"
install -m 0640 -o root -g root "${APP_DIR}/config.yml" "${snapshot_dir}/config.yml"
sqlite3 "$DB_FILE" ".timeout 10000" ".backup '${snapshot_dir}/dujiao.db'"
chown root:root "${snapshot_dir}/dujiao.db"
chmod 0600 "${snapshot_dir}/dujiao.db"
printf '%s  %s\n' "$actual_sha" "$(basename "$artifact")" > "${snapshot_dir}/artifact.sha256"

install -m 0755 -o root -g root "$artifact" "$new_bin"
sync "$new_bin"
mv -f "$new_bin" "$APP_BIN"

rollback() {
  printf 'health check failed; rolling back to previous binary and database\n' >&2
  systemctl stop "$SERVICE" || true
  install -m 0755 -o root -g root "${snapshot_dir}/dujiao-next" "$APP_BIN"
  rm -f "${DB_FILE}-wal" "${DB_FILE}-shm"
  install -m 0640 -o dujiao-next -g dujiao-next "${snapshot_dir}/dujiao.db" "$DB_FILE"
  systemctl start "$SERVICE"
  for _ in $(seq 1 20); do
    if curl --fail --silent --show-error --max-time 2 "$HEALTH_URL" >/dev/null; then
      printf 'rollback completed and previous version is healthy\n' >&2
      return 0
    fi
    sleep 1
  done
  printf 'rollback failed: previous service is not healthy\n' >&2
  return 1
}

if ! systemctl restart "$SERVICE"; then
  rollback
  exit 1
fi

healthy=false
for _ in $(seq 1 30); do
  if curl --fail --silent --show-error --max-time 2 "$HEALTH_URL" >/dev/null; then
    healthy=true
    break
  fi
  sleep 1
done

if [[ "$healthy" != true ]]; then
  journalctl -u "$SERVICE" -n 60 --no-pager >&2 || true
  rollback
  exit 1
fi

printf '%s\n' "$commit" > "${APP_DIR}/DEPLOYED_COMMIT"
printf '%s\n' "$actual_sha" > "${APP_DIR}/DEPLOYED_SHA256"
chmod 0644 "${APP_DIR}/DEPLOYED_COMMIT" "${APP_DIR}/DEPLOYED_SHA256"
rm -f "$artifact"

# Keep the ten most recent pre-deploy snapshots.
find "$BACKUP_ROOT" -mindepth 1 -maxdepth 1 -type d -printf '%T@ %p\n' \
  | sort -rn \
  | tail -n +11 \
  | cut -d' ' -f2- \
  | xargs -r rm -rf --

logger -t zx-shop-deploy "deployed commit=${commit} sha256=${actual_sha}"
printf 'deployment successful: commit=%s sha256=%s\n' "$commit" "$actual_sha"

