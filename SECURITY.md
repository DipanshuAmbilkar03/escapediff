# Security

Report vulnerabilities through [GitHub's private vulnerability reporting](https://github.com/GillMacedo/escapediff/security/advisories/new), not in a public issue.

`escapediff` reads compiler output, runs `git archive` and `go build`, and never uses the network itself or runs the code it builds. Please report privately if a crafted log, repository or archive could make it write outside its temporary directory or do anything else unexpected. A wrong or missed finding is a normal bug.
