const fs = require('fs')

const uploadReleaseAssets = module.exports = async ({ context, github }, tag) => {
  const files = fs.readdirSync('bin').filter(file => file.endsWith('.tar.gz'))

  const release = await github.rest.repos.getReleaseByTag({
    owner: context.repo.owner,
    repo: context.repo.repo,
    tag: tag.replace('refs/tags/', '')
  })

  for (const file of files) {
    console.log(`Uploading bin/${file} to release ${tag.replace('refs/tags/', '')} id=${release.data.id}`)

    await github.rest.repos.uploadReleaseAsset({
      owner: context.repo.owner,
      repo: context.repo.repo,
      release_id: release.data.id,
      name: file,
      data: fs.readFileSync(`bin/${file}`),
    })
  }
}
