# lurkr

Lurkr is an automatic downloading tool in a similar vein to [autodl-irssi](https://github.com/autodl-community/autodl-irssi), [PyWhatAuto](https://github.com/jimrollenhagen/pywhatauto) and
others.

For ease of use, It can be deployed as a single compiled executable file + YAML config.


## Features (Planned)

- [ ] **Per tracker transports** - Each tracker can use its own transport for handling torrents files. This means
you can use a single instance to manage multiple servers or torrent instances for example.- [ ]
- [ ] **Download Filters** Filters can be defined per-source, shared, or both. 
    - [ ] **Name** - Match specific names, eg: "The Simpsons"
    - [ ] **Exclusions** - Never download anything containing the exlusions  
    - [ ] **Min/Max Size**
    - [ ] **Categories** - TV, Movies, Music, Games, Apps, Anime, Other
    - [ ] **Genres** 
    - [ ] **Formats** - mkv, flac
    - [ ] **Resolution** - 720p, 1080p, 4K      
    - [ ] **Seasons**
    - [ ] **Episode**
    - [ ] **Sources** HDTV, WEB-DL, BLURAY, etc.  
    - [ ] **Age** (Release date, Max Pre time)
    - [ ] **Score** (Imdb, etc.)
- [ ] **ZNC Integration**
- [ ] **WebUI** Simple  for configuration & status display
- [ ] **Duplication detections** Optionally only download a single version of any given show, tv, etc.
    - [ ] Duplicate infohash
    - [ ] Duplicate releases
### Transports

- Watch dirs
    - [x] **FileSystem** - This will download torrents into a specific local directory
    - [x] **SFTP** - Upload the torrent to a directory on a remote server
    - [ ] **FTP** - If there is demand, not recommended
    
- Clients
    - [ ] **[Deluge](https://deluge-torrent.org/) RPC** - 
    - [ ] **[Transmission](https://transmissionbt.com/) JSON-RPC**
    - [ ] **[rTorrent](https://github.com/rakshasa/rtorrent) SCGI**
    - [ ] **[qBittorrent](https://www.qbittorrent.org/) Web API**
    - **Not Planned**: uTorrent, Tixati, Vuze, BiglyBT
    
### Tracker Source Support

- **IRC Announces**
    - [x] **BTN**
    - [ ] **PTP**
    - [x] **TL**
    - [x] **RevTT**
    - [x] **RED**
    - [ ] ...many more, PRs welcomed.
    
- **RSS Feeds**
    - [ ] **Generic**
    - [ ] **BLU**
    - [ ] ...many more, PRs welcomed.
    
## Docker

There will be an included docker file you can build yourself as well as an image available, ready to run.

## Seedr

There is a sister application called [seedr](https://github.com/leighmacdonald/seedr) which helps you manage your
torrents for optimizing seeding, using a seedbox or otherwise. There will eventually be some (optional) integrations between the
two applications.