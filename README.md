## RSS Feed Aggregator

Building an aggregator service in Go. It's a backend web server which allows a user to:
- Follow 5 RSS feeds by default (on successful login)
- Follow and unfollow RSS feeds that other users have added
- Fetch all of the latest posts from the RSS feeds they follow

<br>

- [JWT Authentication Service](https://github.com/siddhant-vij/JWT-Authentication-Service) to securely login users on the website.
- [Dynamic Blog Generator](https://github.com/siddhant-vij/Dynamic-Blog-Generator) using Go's template engine to display the RSS feed posts.

<br>

RSS feeds are a way for websites to publish updates to their content. You can use this project to keep up with your favorite blogs, news sites, podcasts, and more!

[Boot.dev](https://www.boot.dev/learn/build-blog-aggregator) Guided Project

<br>

### 🚀 Learning Goals
- How to integrate a Go server with PostgreSQL
- The basics of database migrations
- Long-running service workers
- The complete overview of JWT authentication
- Dynamic blog generator using the backend API

<br>

### 🚀 Improvement Ideas
- Support different options for sorting and filtering posts using query parameters
- Support multiple types of RSS feeds with better logging and error handling (e.g. Atom, JSON, etc.)
- Classify different types of feeds and posts (e.g. blog, podcast, video, etc.)
- Support pagination of the endpoints that can return many items
- Add a CLI client that uses the API to fetch and display posts, maybe it even allows you to read them in your terminal
- Scrape lists of feeds themselves from a third-party site that aggregates feed URLs (e.g. FeedSpot, etc.)
- Add integration tests that use the API to create, read, update, and delete feeds and posts
- Add bookmarking or "liking" to posts
- Create a comprehensive web UI that uses the backend API