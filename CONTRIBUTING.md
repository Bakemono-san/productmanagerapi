# Contributing to Product Manager API

Thank you for your interest in contributing to the Product Manager API! We welcome contributions from everyone and appreciate your help in making this project better.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Environment Setup](#development-environment-setup)
- [How to Contribute](#how-to-contribute)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Commit Message Guidelines](#commit-message-guidelines)
- [Pull Request Process](#pull-request-process)
- [Issue Reporting](#issue-reporting)
- [Documentation](#documentation)
- [Getting Help](#getting-help)

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct. We are committed to providing a welcoming and inspiring community for all.

### Our Standards

* Be respectful and inclusive
* Be collaborative and constructive
* Focus on what is best for the community
* Show empathy towards other community members

## Getting Started

### Prerequisites

Before you begin, ensure you have the following installed:

* **Go 1.16+** (Go 1.23.5 recommended)
* **PostgreSQL** (version 12 or higher)
* **Git**
* **A code editor** (VS Code, GoLand, or similar)

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:

```bash
git clone https://github.com/YOUR_USERNAME/productmanagerapi.git
cd productmanagerapi
```

3. Add the original repository as upstream:

```bash
git remote add upstream https://github.com/Bakemono-san/productmanagerapi.git
```

## Development Environment Setup

### Database Setup

1. Create a PostgreSQL database:

```bash
createdb product_management
```

2. Update the database configuration in `config/dbconfig.go`:

```go
const dsn = "host=localhost user=youruser password=yourpass dbname=product_management port=5432 sslmode=disable"
```

### Install Dependencies

```bash
go mod download
```

### Run the Application

```bash
go run cmd/main.go
```

The server should start on `http://localhost:2002`

### Verify Setup

Test the setup by accessing:
- API: `http://localhost:2002/`
- Swagger docs: `http://localhost:2002/swagger/index.html`

## How to Contribute

### 1. Choose What to Work On

* Check the [Issues](https://github.com/Bakemono-san/productmanagerapi/issues) for open tasks
* Look for issues labeled `good first issue` if you're new to the project
* Feel free to propose new features by creating an issue first

### 2. Create a Branch

```bash
# Update your local main branch
git checkout main
git pull upstream main

# Create a new feature branch
git checkout -b feature/your-feature-name

# Or for bug fixes
git checkout -b fix/issue-description
```

### 3. Make Your Changes

* Write clean, readable code
* Follow the existing code style
* Add tests for new functionality
* Update documentation as needed

### 4. Test Your Changes

```bash
# Run the application
go run cmd/main.go

# Test your changes manually
# Verify database migrations work
# Check that existing functionality still works
```

### 5. Commit and Push

```bash
git add .
git commit -m "feat: add user profile management"
git push origin feature/your-feature-name
```

### 6. Create a Pull Request

* Go to GitHub and create a pull request
* Fill out the pull request template
* Link any related issues

## Coding Standards

### Go Style Guidelines

* Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
* Use `gofmt` to format your code
* Use `golint` to check for style issues
* Follow Go naming conventions:
  * **PascalCase** for exported functions and types
  * **camelCase** for unexported functions and variables
  * **ALL_CAPS** for constants

### Project Structure

When adding new code, follow the existing project structure:

```
├── controllers/     # HTTP request handlers
├── services/        # Business logic layer
├── models/          # Database models
├── routes/          # API route definitions
├── utils/           # Utility functions and middleware
├── config/          # Configuration files
└── types/           # Custom type definitions
```

### Code Quality

* **Functions should be small** and focused on a single responsibility
* **Use meaningful names** for variables, functions, and types
* **Add comments** for exported functions and complex logic
* **Handle errors properly** - don't ignore errors
* **Use dependency injection** where appropriate

### Example Code Style

```go
// Good: Clear function name and documentation
// GetUserByID retrieves a user from the database by their ID.
// Returns an error if the user is not found.
func GetUserByID(id uint) (*models.User, error) {
    var user models.User
    result := config.Db.First(&user, id)
    if result.Error != nil {
        return nil, fmt.Errorf("user not found: %w", result.Error)
    }
    return &user, nil
}

// Bad: Unclear naming and no error handling
func getUser(id uint) *models.User {
    var user models.User
    config.Db.First(&user, id)
    return &user
}
```

## Testing Guidelines

### Manual Testing

When adding new features:

1. **Test the happy path** - ensure normal operation works
2. **Test edge cases** - empty inputs, invalid data, etc.
3. **Test error conditions** - network failures, database errors
4. **Test authentication** - ensure protected endpoints require tokens

### API Testing with cURL

```bash
# Test user registration
curl -X POST http://localhost:2002/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"testpass","role":"user"}'

# Test login
curl -X POST http://localhost:2002/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"testpass"}'

# Test protected endpoint (replace TOKEN with actual token)
curl -X GET http://localhost:2002/products \
  -H "Authorization: Bearer TOKEN"
```

### Integration Testing

* Verify database migrations work correctly
* Test the complete workflow: register → login → create/read/update/delete
* Ensure authentication middleware works properly
* Test CORS functionality

## Commit Message Guidelines

We use [Conventional Commits](https://www.conventionalcommits.org/) for consistent commit messages.

### Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types

* **feat**: A new feature
* **fix**: A bug fix
* **docs**: Documentation only changes
* **style**: Changes that do not affect the meaning of the code
* **refactor**: A code change that neither fixes a bug nor adds a feature
* **perf**: A code change that improves performance
* **test**: Adding missing tests or correcting existing tests
* **chore**: Changes to the build process or auxiliary tools

### Examples

```bash
feat: add user profile update endpoint
fix: resolve JWT token expiration issue
docs: update API documentation for sales endpoints
refactor: improve error handling in product controller
test: add unit tests for authentication service
chore: update dependencies and go.mod
```

### Detailed Example

```
feat(auth): add password reset functionality

- Add password reset request endpoint
- Add password reset confirmation endpoint
- Add email notification service
- Update user model with reset token fields

Closes #123
```

## Pull Request Process

### Before Creating a PR

- [ ] Code follows the project's style guidelines
- [ ] Self-review of the code has been performed
- [ ] Code is commented, particularly in hard-to-understand areas
- [ ] Documentation has been updated if necessary
- [ ] Changes have been tested locally
- [ ] No merge conflicts with the main branch
- [ ] Commit messages follow the guidelines

### PR Description Template

```markdown
## Description
Brief description of the changes made in this PR.

## Type of Change
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update

## How Has This Been Tested?
- [ ] Unit tests
- [ ] Integration tests
- [ ] Manual testing
- [ ] Tested with Postman/cURL

## Screenshots (if applicable)
Add screenshots or API response examples to help reviewers understand the changes.

## Checklist
- [ ] My code follows the style guidelines of this project
- [ ] I have performed a self-review of my own code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings or errors
- [ ] Any dependent changes have been merged and published

## Related Issues
Closes #(issue_number)
```

### Review Process

1. **Automated checks** must pass (if any are set up)
2. **At least one maintainer** must review the PR
3. **Address feedback** from reviewers promptly
4. **Maintainers will merge** approved PRs

## Issue Reporting

### Bug Reports

When reporting bugs, please include:

**Bug Report Template:**

```markdown
## Bug Description
A clear and concise description of what the bug is.

## Steps to Reproduce
1. Go to '...'
2. Click on '....'
3. Scroll down to '....'
4. See error

## Expected Behavior
A clear and concise description of what you expected to happen.

## Actual Behavior
A clear and concise description of what actually happened.

## Screenshots/Logs
If applicable, add screenshots or error logs to help explain the problem.

## Environment
- OS: [e.g. Ubuntu 20.04, macOS 12.0, Windows 10]
- Go version: [e.g. 1.23.5]
- PostgreSQL version: [e.g. 14.2]
- Browser (if applicable): [e.g. Chrome 95.0]

## Additional Context
Add any other context about the problem here.
```

### Feature Requests

**Feature Request Template:**

```markdown
## Feature Description
A clear and concise description of what you want to happen.

## Problem Statement
What problem does this feature solve? Why is it needed?

## Proposed Solution
Describe the solution you'd like to see implemented.

## Alternatives Considered
Describe any alternative solutions or features you've considered.

## Additional Context
Add any other context, screenshots, or examples about the feature request here.

## Implementation Notes
If you have ideas about how this could be implemented, share them here.
```

## Documentation

### API Documentation

* Update Swagger comments when adding new endpoints
* Include request/response examples
* Document all parameters and their types
* Specify authentication requirements

### Code Documentation

* Add comments for all exported functions and types
* Include examples in comments for complex functions
* Document function parameters and return values
* Explain the purpose of complex algorithms

### README Updates

Update the README when you:
* Add new features
* Change installation/setup requirements
* Modify API endpoints
* Add new dependencies

## Getting Help

### Where to Get Help

* **GitHub Issues**: For bugs, feature requests, and general questions
* **GitHub Discussions**: For community discussions and help
* **Code Review**: Ask for help in your pull request

### What to Include When Asking for Help

* Clear description of what you're trying to achieve
* What you've already tried
* Error messages or unexpected behavior
* Relevant code snippets
* Your environment details

### Response Time

* **Issues**: We aim to respond within 48-72 hours
* **Pull Requests**: Initial review within 3-5 days
* **Security Issues**: Contact maintainers directly for faster response

## Recognition

Contributors will be recognized in the following ways:

* Listed in the project's contributors section
* Mentioned in release notes for significant contributions
* Invited to become maintainers for consistent, high-quality contributions

## Questions?

If you have any questions about contributing, feel free to:

* Open an issue with the `question` label
* Start a discussion in GitHub Discussions
* Reach out to the maintainers

Thank you for contributing to Product Manager API! 🚀

---

*This contributing guide is inspired by open source best practices and is subject to updates as the project evolves.*