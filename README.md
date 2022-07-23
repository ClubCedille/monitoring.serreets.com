# monitoring.serreets.com


Site interne du club [serreÉTS](https://serreets.com) permettant de surveiller une serre expérimentale sur le campus de l'ÉTS.

## Dependencies
1. [go](https://go.dev/)
2. [nodejs](https://nodejs.org/en/)


## Development

### Server
You can start the server by running 
```bash
$ go run api/main.go
```
Visit http://localhost:3001/api/hello in your web browser.

### Web
You can run the web app  in development mode with :
```bash
$ npm start --prefix assets
# or with yarn
$ yarn --cwd assets run start
```
Open http://localhost:3000 to view it in the browser.

## Production 

To build the web app for production, run : 
```bash
$ npm run build --prefix assets
# or with yarn 
$ yarn --cwd assets build 
```
It bundles React in production mode to the build folder

# Docker
From the root directory, run this command:
```bash
docker build -t serreets-backend . && docker build --file=assets/Dockerfile -t serreets-frontend --build-arg REACT_APP_URL=http://127.0.0.1:3001/ . && docker-compose -f docker-compose.yml up
```
