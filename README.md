﻿

## Employee-Doc

---

A RESTful API backend system for managing employee document uploads, including KTP, diploma, work contracts, and NPWP. This backend project ensures secure authentication, proper data validation, file handling, and modular architecture using Go and Clean Architecture principles. The system is production-ready and containerized using Docker.

---

## Installation 👨🏻‍💻

This project is built with Go and uses Go 1.20+.

1. Clone Repository
   By use terminal/cmd

```sh
git clone hhttps://github.com/Kevinmajesta/BookManagement.git (BELUM)
```

2. Open Repository
   By use terminal/cmd

```sh
cd bookmanagement
```

2. Check the .env file and configure with your device

3. Enable the PostgreSQL database
   Option you can use :

   - [pgAdmin](https://www.pgadmin.org/)
   - [NaviCat](https://www.navicat.com/en/download/navicat-premium?gad_source=1&gclid=CjwKCAjwmYCzBhA6EiwAxFwfgFWv6YNc_nwrdL5BByjvaEmUNbzD0vvg-tHgv7x6rFyIx-zSdWYQWhoCRP0QAvD_BwE)
   - Or anything else you usualy use

4. Run the command to create the database and migrate it.
   Make sure you have install migrate cli.
   If you dont, install first by

**If you MAC user** 🍏

- First if you dont have [Home Brew](https://brew.sh/)
  Open terminal and copy code below :

```sh
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

Then install migrate cli

```sh
brew install golang-migrate
```

**If you windows user**🪟

- Open PowerShell. You can do this by searching for PowerShell in the start menu.
- Inside PowerShell, type the code below

```sh
iwr -useb get.scoop.sh | iex
```

Then Install use scoop

```sh
scoop install migrate
```

After all, migrate it by

```sh
migrate -path=database/migrations -database "postgresql://upload-karyawan:upload-karyawan@localhost:5432/upload-karyawan?sslmode=disable" -verbose up


```

5. Configure Docker
   **First Install Docker**
   - Windows User[Docker](https://docs.docker.com/desktop/install/windows-install/)
   - Mac User [Docker](https://docs.docker.com/desktop/install/mac-install/)
   - Then compose it by

```sh
docker-compose up -d
```

```sh
https://hub.docker.com/repository/docker/kevinivano/myapp/general
for the collect the project link that has been built and saved to the registry container hub.docker.com
or you can pull with 
docker pull kevinivano/myapp
```

6. Run the application

```sh
go mod tidy
go run cmd/app/main.go
```

---

## ✨ Features

- **Admin Authentication**
  - Login endpoint for administrator
  - JWT-based authorization

- **Employee Document Management**
  - Upload & download KTP, Ijazah, NPWP, Kontrak Kerja
  - CRUD operations for documents
  - Files stored and accessed via local storage

- **Validation**
  - All inputs validated with proper error responses

- **Seeding & Migration**
  - Auto-seed employee data
  - Migration support via golang-migrate

- **Queue Worker & Cache & Goroutine & Channel**
  - File processing managed through Redis queue
  - Cache using redis
  - Goroutine email sent and upload job
  - channel for RunInternalNotificationWorker


---

## 🛠️ Tech Stack

- **Language**: Go (Golang)
- **Framework**: Gin
- **Database**: PostgreSQL
- **Communication**: RESTful API 
- **Authentication**: JWT
- **Containerization**: Docker
- **Architecture**: Microservices Clean Architecture
- **Cache & Queue**: Redis



## Development

This project app develope by 1 people
| Name | Github |
| ------ | ------ |
| Kevin | https://github.com/Kevinmajesta |


By using github for development for staging and production.


