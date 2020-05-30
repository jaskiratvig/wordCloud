# Word Cloud App

This project uses the Spotify API and Genius API to create a word cloud based on the lyrics of a given Spotify user's top 20 songs.

## Dependancies

The following dependancies are required to run the project:

* Create a Spotify OAuth bearer token:
  * Visit ``` https://developer.spotify.com/console/get-current-user-top-artists-and-tracks/ ```
  * Under OAuth token, click the green box that says "Get Token"
  * Under Required scopes for this endpoint, tick "user-top-read" and click "Request Token"
  * Copy the generated OAuth token and set the "spotifyBearerToken" variable in main.go
* Create a Genius OAuth bearer token:
  * Visit ``` https://api.genius.com/oauth/authorize ``` and create an account
  * Visit ``` https://genius.com/api-clients ``` and create an API client
  * Copy the generated access token and set the "geniusBearerToken" variable in main.go

## Run Project

To run this project, navigate to the wordCloud root directory and run ``` go run main.go ```. The word cloud should populate under output.png. Happy word-clouding!

## Acknowledgements

All code in the cloudCreation folder was pulled from ``` https://github.com/psykhi/wordclouds ```
