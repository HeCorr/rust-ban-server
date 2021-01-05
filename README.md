# rust-ban-server
Fully-featured ban server compliant with the new Rust (game) server's [Centralized Banning](https://wiki.facepunch.com/rust/centralized-banning) feature written in Go + GORM/SQLite.

---

#### Usage

<details>
<summary>Running from the binaries</summary>

1. Download the [latest release](https://github.com/HeCorr/rust-ban-server/releases/latest) using a compatible binary for your system

2. Execute it: `./rust-ban-server`.
Available flags:
    - `-l` API listen address (default: `:4000`)
    - `-q` Quiet mode, omits HTTP log output.

</details>

<details>
<summary>Building from source</summary>

**WIP**
    
</details>


#### Available Endpoints

- `GET /api/status` - For checking if the API is alive
- `GET /api/rustBans/<SteamID64>` - For checking for banned SteamIDs
(mostly used by the game server)
- `POST /api/rustBans` - For adding account bans
- `DELETE /api/rustBans/<SteamID64>` - For removing account bans


#### TODO
- Secure the `POST` and `DELETE` endpoints with a token or access key;
- Actually test the API on a local Rust server;
- HTTPS support;
- Include Go test files;
- Database importing and exporting;
- Create endpoint that returns ban count;
- Create endpoint that returns all bans (might wanna implement pagination tho);

#### Spec ([subject to change](https://youtu.be/YOEd19K9WZA?t=158))
`GET /api/status` shall always return `{ "status": "ok" }` with status code `200`.

`GET /api/rustBans/<SteamID64>` first validates the provided SteamID, returning the status code `400` and the JSON-encoded error message `Invalid SteamID64.` if it's not valid.

Otherwise, it returns the JSON data as specified in the Rust wiki with status code `200` if the SteamID has been found.

If it wasn't, it returns `404` with the JSON-encoded error `SteamID64 not found.`.

In case of internal errors, the API returns `500` and the JSON-encoded error.

`POST /api/rustBans` requires a JSON body to be sent through, using the same format as described in the Rust wiki:
```json
{
    "steamId": "76561198060722078",
    "reason": "Too handsome",
    "expiryDate": 1609698084
}
```
The provided SteamID64 is also validated, returning the status code `400` and an error message if the validation fails.

If the request was successful, the API returns the JSON-encoded message `SteamID64 banned.` with the status code `201`.

If the SteamID already exists, the status code `209` is returned with the JSON-encoded error `SteamID64 already banned.`.

In case of internal errors, it returns `500` and the JSON-encoded error.

`DELETE /api/rustBans/<steamID64>` also validates the provided SteamID64, returning the status code `400` and a JSON-encoded error if validation fails.

If it's valid, the status code `200` is returned with the JSON-encoded message `SteamID64 unbanned.` if the SteamID has been removed successfully.

If the SteamID doesn't exist, the API returns `404` with the JSON-encoded error `SteamID64 not banned.`.

In case of internal errors, it also returns `500` and the JSON-encoded error.