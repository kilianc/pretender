const fs = require('fs')

const updateCodeCoverageComment = module.exports = async ({ context, github }) => {
  const comments = await github.rest.issues.listComments({
    owner: context.repo.owner,
    repo: context.repo.repo,
    issue_number: context.issue.number,
    per_page: 100
  })

  const coverageComment = comments.data.find((comment) => {
    return comment.body.startsWith('<!-- coverage -->')
  }) || {}

  const coverageText = fs.readFileSync('cover.txt', 'utf8').split('\n').slice(0, -1)
  const coverageTextSummary = coverageText[coverageText.length-1].split('\t').pop()

  const commentBody = [
    '<!-- coverage -->',
    `### Code Coverage Report ${process.env.REVISION}`,
    '```',
    `Total: ${coverageTextSummary}`,
    '```',
    '<details>',
    '<summary>Full coverage report</summary>',
    '',
    '```',
    ...coverageText,
    '```',
    '</details>',
  ].join('\n')

  const upsertCommentOptions = {
    owner: context.repo.owner,
    repo: context.repo.repo,
    issue_number: context.issue.number,
    comment_id: coverageComment.id,
    body: commentBody
  }

  if (coverageComment.id) {
    await github.rest.issues.updateComment(upsertCommentOptions)
  } else {
    await github.rest.issues.createComment(upsertCommentOptions)
  }
}
