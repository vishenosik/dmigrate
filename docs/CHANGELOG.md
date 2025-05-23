# Changelog

All notable changes to this project will be documented in this file.

## [0.0.1] - 2025-04-06

### 🚀 Features

- *(migrate)* Add initial files.
- *(client / queries)* Correct versions schema. Add tests. Update deps.
- *(logging)* Add logger init & functional options. Add std logger to handle nil loggers passed to With func.
- *(client)* Route dgraph client init through connect(). This is meant to add abstraction to users.
- *(connect)* Add cancel function to dgraph connection.

### 🐛 Bug Fixes

- *(collections)* Fix Filter() - by function called return of false. Should be continued.

### 🚜 Refactor

- *(migrations)* Refactor  migrationsToApply().
- *(version queries)* Remove naming parameter from version retrieval. Add type search to store less fields.

### 📚 Documentation

- *(changelog)* Add changelog workflow. Add docs.
- *(license)* Move LICENSE to docs.
- *(readme)* Add LICENSE link to README.
- *(readme)* Add minimal instructions to readme.
- *(migrate)* Add minimal comments to exported functions.

### 🧪 Testing

- *(task scripts)* Add taskfile. Add testing script.
- *(docker tests)* Add more time to sleep waiting for docker.
- *(client ctx)* Make a suite with client conn & context cancelation.
- *(github actions)* Add podman installation to github actions.
- *(github actions)* Replace podman with docker in github actions.
- *(github actions)* Cleanup github actions.
- *(migrate)* Add full migration test. Test includes file reading.

### ⚙️ Miscellaneous Tasks

- *(args namings)* Fix arguments order in dgraph queries. Fix namings.
- *(git)* Update gitignore. Add obsidian files, data storage ignore.
- *(coverfile)* Move cover.out creatin to /test dir.
- *(podman integration)* Fix testing script. Add variables & requirements to Taskfile.

<!-- generated by git-cliff -->
