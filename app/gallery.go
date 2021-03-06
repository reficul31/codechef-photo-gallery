package app

/*
The gallery api
*/

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func AlbumHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var out []byte // json output

	switch r.Method {
	case "GET":
		albums := &Albums{}
		if err = DB.db.Where("user_id = ?", currUser.ID).Find(&albums).Error; err != nil {
			write_response(err, w, false, "Couldn't fetch the albums")
			return
		}

		out, err = json.Marshal(albums)
		if err != nil {
			write_response(err, w, false, "Internal Server Error")
			return
		}

	case "POST":
		album := Album{}
		err = json.NewDecoder(r.Body).Decode(&album)
		if err != nil {
			write_response(err, w, false, "Internal Server Error")
			return
		}

		tx := DB.db.Begin()
		album.UserID = uint(currUser.ID)

		if err = tx.Create(&album).Error; err != nil {
			tx.Rollback()
			write_response(err, w, false, "Can't add album")
			return
		}
		tx.Commit()

		out = []byte("Album Added!")

	case "PUT":
		album := Album{}

		err := json.NewDecoder(r.Body).Decode(&album)
		if err != nil {
			write_response(err, w, false, "Internal Server Error")
			return
		}

		tx := DB.db.Begin()
		if err = tx.Model(&album).Updates(map[string]interface{}{"name": album.Name, "privacy": album.Privacy, "description": album.Description}).Error; err != nil {
			tx.Rollback()
			write_response(err, w, false, "Can't update user")
			return
		}
		tx.Commit()

		out = []byte("Album Updated!")

	case "DELETE":
		album := Album{}
		err = json.NewDecoder(r.Body).Decode(&album)
		if err != nil {
			write_response(err, w, false, "Internal Server Error")
			return
		}

		photos := Photos{}
		if err = DB.db.Where("album_id = ?", album.ID).Find(&photos).Error; err != nil {
			write_response(err, w, false, "Couldn't fetch the album photos")
			return
		}

		tx := DB.db.Begin()
		for _, photo := range photos {
			if err = deleteFile(photo.Name); err != nil {
				tx.Rollback()
				write_response(err, w, false, "Can't delete album.")
			}
			if err = tx.Delete(&photo).Error; err != nil {
				tx.Rollback()
				write_response(err, w, false, "Can't delete album.")
			}
		}
		if err = tx.Delete(&album).Error; err != nil {
			tx.Rollback()
			write_response(err, w, false, "Can't delete album.")
		}
		tx.Commit()

		out = []byte("Album Deleted!")
	}

	write_response(nil, w, true, string(out))
	return
}

func PhotoHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var out []byte // json output

	switch r.Method {
	case "GET":
		albumId, err := strconv.ParseInt(r.URL.Query().Get("albumId"), 10, 16)
		if err != nil {
			write_response(err, w, false, "No album specified for photos")
			return
		}
		photos := &Photos{}
		if err = DB.db.Where("album_id = ?", albumId).Find(&photos).Error; err != nil {
			write_response(err, w, false, "Couldn't fetch the photos")
			return
		}

		out, err = json.Marshal(photos)
		if err != nil {
			write_response(err, w, false, "Internal Server Error")
			return
		}

	case "POST":
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("file")
		if err != nil {
			write_response(err, w, false, "Internal Server Error")
			return
		}

		handler.Filename = RandStringRunes(20)

		mimeType := handler.Header.Get("Content-Type")
		switch mimeType {
		case "image/png":
		    err = saveFile(w, file, handler)
		default:
		    write_response(err, w, false, "The format file is not valid. Please upload only png images.")
		    return
		}

		if err != nil {
			write_response(err, w, false, "Internal Server Error")
			return
		}

		description := r.Form["description"][0]
		privacy, err := strconv.Atoi(r.Form["privacy"][0])
		if err != nil {
			write_response(err, w, false, "Internal Server Error")
			return
		}

		albumId, err := strconv.Atoi(r.Form["albumId"][0])
		if err != nil {
			write_response(err, w, false, "Internal Server Error")
			return
		}

		photo := Photo{
			Name: handler.Filename,
			Description: description,
			Privacy: privacy,
			AlbumID: uint(albumId),
			Likes: 0,
		}

		tx := DB.db.Begin()

		if err = tx.Create(&photo).Error; err != nil {
			tx.Rollback()
			write_response(err, w, false, "Can't add photo")
			return
		}
		tx.Commit()

		out = []byte("Photo Added!")

	case "DELETE":
		photo := Photo{}
		err = json.NewDecoder(r.Body).Decode(&photo)
		if err != nil {
			write_response(err, w, false, "Internal Server Error")
			return
		}

		filename := photo.Name
		tx := DB.db.Begin()
		if err = tx.Delete(&photo).Error; err != nil {
			tx.Rollback()
			write_response(err, w, false, "Can't delete photo.")
		}
		tx.Commit()
		err = deleteFile(filename)
		out = []byte("Photo Deleted!")
	}

	write_response(nil, w, true, string(out))
	return
}

func UserHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var out []byte // json output
	defer r.Body.Close()

	switch r.Method {
	case "GET":
		user := User{}
		if DB.db.Select("id, name, gender, email").Where("id=?", currUser.ID).First(&user).RecordNotFound() {
			write_response(err, w, false, "Couldn't find user")
			return
		}

		rows, err := DB.db.Raw("select count(distinct albums.id), count(photos.id) from photos, albums where photos.album_id = albums.id and albums.id in (select id from albums where user_id=?)", currUser.ID).Rows()
		if err != nil {
			write_response(err, w, false, "Internal Server Error")
			return
		}
		var albums, photos int
		rows.Next()
		rows.Scan(&albums, &photos)
		rows.Close()

		respUser := GetUserStruct{
			Name: user.Name,
			Email: user.Email,
			Gender: user.Gender,
			Albums: albums,
			Photos: photos,
		}


		out, err = json.Marshal(respUser)
		if err != nil {
			write_response(err, w, false, "Internal Server Error")
			return
		}
	case "PUT":
		respUser := User{}
		err = json.NewDecoder(r.Body).Decode(&respUser)

		if err != nil {
			write_response(err, w, false, "Internal Server Error")
			return
		}

		updatedUser := User{}
		if DB.db.Where("id=?", currUser.ID).First(&updatedUser).RecordNotFound() {
			write_response(err, w, false, "Couldn't find user")
			return
		}

		updatedUser.Name = respUser.Name
		updatedUser.Gender = respUser.Gender
		updatedUser.Email = respUser.Email

		tx := DB.db.Begin()
		if err = tx.Save(&updatedUser).Error; err != nil {
			tx.Rollback()
			write_response(err, w, false, "Can't update user")
			return
		}
		tx.Commit()
		out = []byte("User Updated!")
	case "DELETE":
		user := User{}
		if DB.db.Select("id, name, gender, email").Where("id=?", currUser.ID).First(&user).RecordNotFound() {
			write_response(err, w, false, "Couldn't find user")
			return
		}

		albums := Albums{}
		if err := DB.db.Where("user_id=?", user.ID).Find(&albums).Error; err != nil {
			write_response(err, w, false, "Couldn't find user's albums")
			return
		}

		tx := DB.db.Begin()
		for _, album := range albums {
			photos := Photos{}
			if err := DB.db.Where("album_id=?", album.ID).Find(&photos).Error; err != nil {
				write_response(err, w, false, "Couldn't find user's albums")
				return
			}
			for _, photo := range photos {
				if err = deleteFile(photo.Name); err != nil {
					tx.Rollback()
					write_response(err, w, false, "Can't delete user.")
				}
				if err = tx.Delete(&photo).Error; err != nil {
					tx.Rollback()
					write_response(err, w, false, "Can't delete user.")
				}
			}
			if err = tx.Delete(&album).Error; err != nil {
				tx.Rollback()
				write_response(err, w, false, "Can't delete user.")
			}
		}

		if err = tx.Delete(&user).Error; err != nil {
			tx.Rollback()
			write_response(err, w, false, "Can't delete user.")
		}

		tx.Commit()
		out = []byte("User Deleted!")
	}

	write_response(nil, w, true, string(out))
	return
}

func FetchPhoto(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	photoId, err := strconv.ParseInt(r.URL.Query().Get("photoId"), 10, 16)
	if err != nil {
		write_response(err, w, false, "Unable to read query parameter")
		return
	}

	photo := Photo{}
	if DB.db.Where("id=?", photoId).First(&photo).RecordNotFound() {
		write_response(err, w, false, "Can't find the photo.")
		return
	}

	var userId int
	if photo.Privacy == 0 {
		row, err := DB.db.Raw("select users.id from users, albums where users.id = albums.user_id and albums.id = (select album_id from photos where photos.id = ?)", photo.ID).Rows()
		if err != nil {
			write_response(err, w, false, "Internal Server Error.")
			return
		}
		row.Next()
		row.Scan(&userId)
		row.Close()
		if userId != currUser.ID {
			write_response(err, w, false, "This photo seems to be private")
			return
		}
	}

	out, err := json.Marshal(photo)
	if err != nil {
		write_response(err, w, false, "Internal Server Error.")
		return
	}

	write_response(nil, w, true, string(out))
	return
}

func FetchAlbum(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	albumId, err := strconv.ParseInt(r.URL.Query().Get("albumId"), 10, 16)
	if err != nil {
		write_response(err, w, false, "Unable to read query parameter")
		return
	}

	album := Album{}
	if DB.db.Where("id=?", albumId).First(&album).RecordNotFound() {
		write_response(err, w, false, "Can't find the album.")
		return
	}

	if album.Privacy == 0 && album.UserID != uint(currUser.ID) {
		write_response(err, w, false, "This seems to be a private album.")
		return
	}

	photos := Photos{}
	if err = DB.db.Where("album_id=?", album.ID).Find(&photos).Error; err != nil {
		write_response(err, w, false, "Can't find the photos of the album.")
		return
	}

	album.Photos = photos

	out, err := json.Marshal(album)
	if err != nil {
		write_response(err, w, false, "Internal Server Error.")
		return
	}

	write_response(nil, w, true, string(out))
	return
}