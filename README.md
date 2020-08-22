# Chipper :: go-social-network

<img src="https://i.imgur.com/le0Ghuc.png" />

# Introduction
This is a mini social network created with golang, gorilla/mux and gorm. I built this to get more familiar with the Go language. 
I also decided to use go template just to get familiar with them. I'll probably go with rest apis for future projects.
I intend to do some refactoring later on, but the functionalities are currently working. Feel free to fork it or make some suggestions to the current implementation.
Inspiration gotten from <a href="/https://github.com/yTakkar/Go-Mini-Social-Network">yTakkar's</a> .
The website is live at https://chipper.daredev.xyz

# Requirements
Go - from v1.2

PostgreSQL - from v10

# Setup
- Before setup, you must have PostgreSQL on your machine. Then clone this repository

    > https://github.com/dsa0x/go-social-network.git

- cd into the project directory

   > cd go-social-network

- Set environment variables

  - A sample is presented in the repository.
  
  > cp .env.example .env

- Run application

  > go run main.go


# Usage

Endpoint | Usage
------------ | -------------
/ | The Home page. It also lists the posts of all users
/signup | Use this route to sign up a new user
/login | Authenticate user to access protected endpoints
/user/{id} | The User profile. Also showing the number of posts, followers, and followings of the user
/logout | Logout the user
/user/follow | To follow a user
/user/unfollow | To unfollow a user
/users | Displays a list of all users
/post/create | To create a new post
/post/{id}/update | status.NotYetImplemented




