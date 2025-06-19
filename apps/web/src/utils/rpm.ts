//https://models.readyplayer.me/685470d6259f6d6e519467ee.glb
export function getPreviewAvatar(avatarUrl: string) {
  if (!avatarUrl) return

  const [url, avatarId] = avatarUrl.slice(8).split('/')
  const urlSplit = url.split('.')
  urlSplit[0] = "api"

  return `https://${urlSplit.join(".")}/v2/avatars/${avatarId}?preview=true`
}
