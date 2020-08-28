# FGL Database

This is the server for the FGL Database, a terminal based communication platform for the Future Gadget Lab. This started as a joke, but I enjoyed developing it, and my friends enjoyed using it.

You can build the server and [client](https://github.com/trebabcock/fgl-client) to create your own server, based on your own needs.

# Features

### Current

- User accounts (register/login)
- Announcements
- Lab Reports
- Gadget Reports
- Member Reports
- Interactive shell

### In Development  

- Group chat
- Private messages
- Multiplayer games

# Building

### Requirements

- Go 1.14+
- PostgreSQL
- Windows, macOS, or Linux

### Compiling

```
git clone https://github.com/trebabcock/fgl-backend.git
cd fgl-backend
go build
```

### Configuration
Before you are able to run the server, you must have a PostgreSQL database running, and you must edit `db/dbConfig.json` with your database's information.

You may also choose a port in `config/config.json`, or leave it as `2814`. More configuration options are coming soon.
