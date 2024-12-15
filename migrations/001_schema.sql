CREATE TABLE users(
    id          VARCHAR(255) NOT NULL PRIMARY KEY,
    username    VARCHAR(255),
    created     DATE
   
);
CREATE TABLE sessions(
    username    VARCHAR(255),
    session     VARCHAR(255) 
);

CREATE TABLE url(
    longurl     VARCHAR(255),
    shorturl    VARCHAR(255)
);
