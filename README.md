# Chess-TS V4
Author: jsbento

## Frontend
- Typescript
- React
- Vite

### Adding Packages
After adding new packages, the containers must be stopped and restarted with `docker compose up --build`

## Backend
- Golang
- Postgres

### Chess Engine
The engine is a Golang port of the Vice chess engine developed by [Bluefever Software](https://www.youtube.com/@BlueFeverSoft). The engine development series can be found [here](https://www.youtube.com/playlist?list=PLZ1QII7yudbc-Ky058TEaOstZHVbT-2hg). The data used by the opening book can be found [here](https://github.com/lichess-org/chess-openings) as is used under a CC0 v1 License.
