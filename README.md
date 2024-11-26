# **BellowAverageMoviesDatabase**

This project demonstrates the creation and querying of a database built using SQLite and populated with movie-related data extracted from multiple CSV files. The project is written in Go (`Golang`) and involves tasks such as handling CSV imports, defining relational database schemas, and performing complex queries. It also covers debugging, error handling, and writing query results to files.

---

## **Features**

1. **Database Creation and Population**:
   - Creates tables for actors, directors, genres, movies, and roles using SQLite.
   - Imports data from multiple CSV files into the database, ensuring data integrity.

2. **Query Execution**:
   - Performs complex queries, such as:
     - Listing the top movies in a specific genre.
     - Displaying directors and their associated genres.
     - Showing actors and their roles in movies.

3. **Result Export**:
   - Stores query results in separate output text files for future reference.


## **Database Schema**

The database is composed of six tables. Below is the data dictionary, detailing the structure and purpose of each table:

| **Table Name**       | **Description**                                                                 | **Columns**                                                                                               |
|-----------------------|---------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------|
| **actors**           | Stores information about actors.                                                | `actor_id` (INTEGER, Primary Key): Unique identifier for each actor.                                     |
|                       |                                                                                 | `actor_name` (TEXT): Name of the actor.                                                                  |
| **directors**        | Stores information about directors.                                             | `director_id` (INTEGER, Primary Key): Unique identifier for each director.                               |
|                       |                                                                                 | `director_name` (TEXT): Name of the director.                                                            |
| **directors_genres** | Maps directors to their associated genres.                                       | `director_id` (INTEGER, Foreign Key): Links to `directors.director_id`.                                  |
|                       |                                                                                 | `genre` (TEXT): The genre associated with the director.                                                  |
| **movies**           | Stores information about movies.                                                | `movie_id` (INTEGER, Primary Key): Unique identifier for each movie.                                     |
|                       |                                                                                 | `movie_name` (TEXT): Name of the movie.                                                                  |
|                       |                                                                                 | `movie_year` (INTEGER): Year the movie was released.                                                     |
|                       |                                                                                 | `movie_rank` (REAL): The rating or ranking of the movie.                                                 |
| **movies_genres**    | Maps movies to their associated genres.                                          | `movie_id` (INTEGER, Foreign Key): Links to `movies.movie_id`.                                           |
|                       |                                                                                 | `genre` (TEXT): The genre associated with the movie.                                                     |
| **roles**            | Maps actors to their roles in specific movies.                                   | `role_id` (INTEGER, Primary Key): Unique identifier for each role.                                       |
|                       |                                                                                 | `actor_id` (INTEGER, Foreign Key): Links to `actors.actor_id`.                                           |
|                       |                                                                                 | `movie_id` (INTEGER, Foreign Key): Links to `movies.movie_id`.                                           |
|                       |                                                                                 | `role_name` (TEXT): The name or title of the role the actor played.                                      |

---

## **Programming Highlights**

### **CSV Import**
- Handled with the `encoding/csv` package.
- Skips malformed rows or rows with fewer fields than expected.
- Uses prepared SQL statements for efficient batch inserts.

### **Queries**
- **Top Movies in a Genre**:
  Fetches the top-ranked movies in a specific genre.
- **Directors and Their Genres**:
  Lists directors and their associated genres.
- **Actors and Their Roles**:
  Shows actors, the roles they played, and the movies they appeared in.

### **Exporting Query Results**
- Each query writes its results to a separate text file.
- Results are formatted for readability.

---

## **Challenges and Solutions**

### **1. Handling Malformed CSV Rows**
- **Issue**: Some rows in the CSV files have fewer fields than expected.
- **Solution**: Added logic to skip such rows silently.

### **2. Slow Compilation with `modernc.org/sqlite`**
- **Issue**: Importing `modernc.org/sqlite` slows down compilation due to its pure Go SQLite implementation.
- **Solution**: Modularized the code to isolate database operations in separate programs (`create_db` and `query_db`).

### **3. No Results from Queries**
- **Issue**: Some queries initially returned no results due to missing or incorrect data relationships.
- **Solution**:
  - Debugged SQL queries directly in SQLite.
  - Ensured that data relationships were correctly established during import.
  - Added debug logs to verify query execution.

### **4. Writing Results to Files**
- **Issue**: Needed to store query results persistently for analysis.
- **Solution**: Implemented a utility function to write query results to files in a human-readable format.

---

## **Sample Outputs**

### **Top Movies in Genre 'Action' (`top_movies_action.txt`)**:
```
Top 10 Movies in Genre 'Action':
Name: Avengers: Endgame, Year: 2019, Rank: 8.4
Name: Mad Max: Fury Road, Year: 2015, Rank: 8.1
...
```

### **Directors and Their Genres (`directors_genres.txt`)**:
```
Directors and their Genres:
Director: Steven Spielberg, Genre: Drama
Director: Quentin Tarantino, Genre: Action
...
```

### **Actors and Their Roles (`actors_roles.txt`)**:
```
Actors and their Roles:
Actor: Robert Downey Jr., Role: Iron Man, Movie: Avengers
Actor: Scarlett Johansson, Role: Black Widow, Movie: Avengers
...
```

---