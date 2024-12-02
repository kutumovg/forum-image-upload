# Forum Image Upload Project
In forum image upload, registered users have the possibility to create a post containing an image as well as text.

## Objectives

<ul>
    <li>Enable communication between users via posts and comments.</li>
    <li>Allow categorization of posts.</li>
    <li>Support liking and disliking of posts and comments.</li>
    <li>Implement filtering of posts based on categories, created posts, or liked posts.</li>
    <li>Enable sign in or sign up.</li>
</ul>

## User Authentication
Users must register and log in to interact with the forum:

<ul>
    <li>Registration requires an email, username, and password.</li>
    <li>Returns an error if the email is already taken.</li>
    <li>Passwords should be encrypted before storage (Bonus).</li>
    <li>Login session uses cookies to maintain a single session per user, with a defined expiration.</li>
</ul>

## User Interaction

Post and Comment Creation:
<ul> 
    <li>Registered users can create posts and comments. Posts can be associated with categories. Images can be upload to posts</li>
</ul>

Likes and Dislikes: 
<ul>
    <li>Registered users can like or dislike posts and comments. Totals are visible to all users.</li>
</ul>

## Filtering
A filtering mechanism allows users to filter posts by:

<ul>
    <li>Categories (Autobiography, Comedy, Science Fiction, Fantasy, Mystery, Other).</li>
    <li>Created posts (specific to logged-in users).</li>
    <li>Liked posts (specific to logged-in users).</li>
</ul>

## Docker Integration

This project is containerized with Docker:
<ul>
    <li>docker build -t forum .</li>
    <li>docker run -p 8080:8080 forum</li>
</ul>

## Usage
<ul>
    <li>Make sure you have Go installed on your system.</li>
    <li>Clone this repository using your terminal:</li>

```
git clone git@git.01.alem.school:gkutumov/forum-image-upload.git
```

<li>Navigate to the project directory:</li>

```
cd forum
```  
<li>Start the project:</li>

```
go run .
``` 
<li>Set Up Database:</li> 
Use SQLite to initialize the database tables for users, posts, comments, and categories.

<li>Run the Application:</li>
Access the forum locally to start posting and interacting!

</ul>

## Authors
<ul>
    <li>gkutumov</li>
    <li>ntoksano</li>
</ul>

Thank you for using our Forum project! If you have any questions or feedback, feel free to reach out. Happy coding! 🚀