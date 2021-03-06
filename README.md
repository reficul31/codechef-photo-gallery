# CodeChef Photo Gallery
A photo gallery app made for codechef recruitments 2018

## How to Run the Code?
This section describes how the app can be viewed in action.

* Step 1: Download and setup GoLang  
Install GoLang. Official documentation installation is prefect and will work seamlessly. Please be wary to setup the directory structure as specified in the docs.

* Step 2: Download the App to your local machine  
Once golang is installed run the following command  
```
go get github.com/reficul31/codechef-photo-gallery
```
This would fetch the entire application to the following path ```$GOPATH/src/github.com/reficul31/codechef-photo-gallery```

* Step 3: Running the Server  
Once the application is downloaded, run the following commands in sequence to run the server. 
```
go run main.go
```
Once the application is running, go to the following [link](http://localhost:8080).
The dummy account's username and password is  
```
Username: test@test.com
Password: password
```

## API Table
This section describes the API made for and utilized by this app.

### Albums
| HTTP Verb | CRUD   | PARAMS(JSON) | RETURNS(JSON)  | ENDPOINT   |
|-----------|--------|--------------|----------------| -----------|
| GET       | Read   |              | List of Albums | /album     |
| PUT       | Update | Album ID     |                | /album     |
| POST      | Create | Album Object |                | /album     |
| DELETE    | Delete | Album ID     |                | /album     |

### Photo 
| HTTP Verb | CRUD   | PARAMS(JSON) | RETURNS(JSON)  | ENDPOINT   |
|-----------|--------|--------------|----------------| -----------|
| GET       | Read   |              | List of Photos | /photo     |
| PUT       | Update | Photo ID     |                | /photo     |
| POST      | Create | Photo Object |                | /photo     |
| DELETE    | Delete | Photo ID     |                | /photo     |

### User
| HTTP Verb | CRUD   | PARAMS(JSON) | RETURNS(JSON)  | ENDPOINT |
|-----------|--------|--------------|----------------| ---------|
| GET       | Read   |              | User Details   | /user    |
| PUT       | Update | User ID      |                | /user    |
| POST      | Create | User Object  |                | /register|
| DELETE    | Delete | User ID      |                | /user    |

## TODOS
This section contains all the landmarks that have to be hit in order to complete the challenge successfully
### Album  
- [x] An album should have a description (you should implement this, should not be a required field)
- [x] An album should have a cover photo
- [x] An album should be editable.
- [x] Should show date and time of creation
- [x] An album should be removable (delete)
- [x] There should be an albums page.
- [x] An album can have many photos (limit to 1000)
- [x] An album can be public/private/only people with the url can view (privacy settings)
- [ ] An album can be liked by a logged in user.
- [ ] Bonus: geo location
### Photo
- [x] A photo should have a description (you should implement this, should not be a required field)
- [x] A photo can only be uploaded to an album
- [x] Should show date and time of creation
- [x] A photo should be removable (delete)
- [x] A photo can be public/private/only people with the url can view (privacy setting)
- [ ] A photo can be liked by a logged in user.
- [ ] Bonus: geo location
### User
- [x] User module should contain
    - username, First name, last name, email, gender, profile picture, password
    - all editable, other than username
- [x] There should be a user profile page
### Bonus
- [x] Token based authentication on all APIs, instead of session based.
- [x] Implementation of Forgot and Reset Password end points. 
- [ ] Social login via Facebook, Google+, Twitter etc. (OAuth Login)
- [ ] Using of better software engineering principles in the code. 
