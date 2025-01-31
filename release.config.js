module.exports = {
    branches: ['main'],
    repositoryUrl: 'git@github.com:xolinar/kademlia-go.git',
    "plugins": [
        [
            "@semantic-release/commit-analyzer",
            {
                "preset": "conventionalcommits"
            }
        ],
        [
            "@semantic-release/release-notes-generator",
            {
                "preset": "conventionalcommits"
            }
        ],
        [
            "@semantic-release/changelog",
            {
                "changelogFile": "CHANGELOG.md",
                "changelogTitle": "# Changelog\n\n All release notes."
            }
        ],
        [
            "@semantic-release/github",
            {
            
            }
        ],
        [
            "@semantic-release/git",
            {
                "assets": ["CHANGELOG.md", "package.json"],
                "message": "chore(release): ${nextRelease.version} [skip ci]\n\n${nextRelease.notes}",
            }
        ]
    ]
}
  