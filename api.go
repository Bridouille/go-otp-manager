package main

import (
	"database/sql"
	"github.com/go-martini/martini"
	gotp "github.com/hgfischer/go-otp"
	"github.com/martini-contrib/render"
	"net/http"
	"strconv"
	"strings"
)

type Retval struct {
	Id   int
	Name string
	Time int
	Totp string
}

type Otp struct {
	Id     int
	Name   string
	Key    string
	Time   int
	Digits int
}

func getOtps(r render.Render, db *sql.DB) {
	rows, err := db.Query("SELECT * FROM otp")
	if err != nil {
		r.JSON(http.StatusInternalServerError, map[string]string{"msg": err.Error()})
		return
	}

	otps := make([]Retval, 0)
	for rows.Next() {
		var otp Retval
		var key string
		var digits int

		err = rows.Scan(&otp.Id, &otp.Name, &key, &otp.Time, &digits)
		if err != nil {
			r.JSON(http.StatusInternalServerError, map[string]string{"msg": err.Error()})
			return
		}
		totp := &gotp.TOTP{Secret: strings.ToUpper(key), IsBase32Secret: true, Length: uint8(digits), Period: uint8(otp.Time)}
		otp.Totp = totp.Get()
		otps = append(otps, otp)
	}

	r.JSON(http.StatusOK, otps)
}

func getOtp(r render.Render, params martini.Params, db *sql.DB) {
	id := params["id"]

	var otp Retval
	var key string
	var digits int
	err := db.QueryRow("SELECT * FROM otp WHERE id = ?", id).Scan(&otp.Id, &otp.Name, &key, &otp.Time, &digits)
	if err == sql.ErrNoRows {
		r.JSON(http.StatusNotFound, map[string]string{"msg": "No such Otp"})
		return
	} else if err != nil {
		r.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
		return
	}

	totp := &gotp.TOTP{Secret: strings.ToUpper(key), IsBase32Secret: true, Length: uint8(digits), Period: uint8(otp.Time)}
	otp.Totp = totp.Get()
	r.JSON(http.StatusOK, otp)
}

func postOtp(r render.Render, req *http.Request, db *sql.DB) {
	var otp Otp
	var err error

	otp.Name = req.PostFormValue("name")
	if otp.Name == "" {
		r.JSON(http.StatusBadRequest, map[string]string{"msg": "You need to provide a valid account name"})
		return
	}
	otp.Key = req.PostFormValue("key")
	if len(otp.Key) != 32 {
		r.JSON(http.StatusBadRequest, map[string]string{"msg": "The secret key must be 32 bytes length"})
		return
	}
	otp.Time, err = strconv.Atoi(req.PostFormValue("time"))
	if err != nil {
		r.JSON(http.StatusBadRequest, map[string]string{"msg": "You need to provide a valid number as time"})
		return
	}
	otp.Digits, err = strconv.Atoi(req.PostFormValue("digits"))
	if err != nil {
		r.JSON(http.StatusBadRequest, map[string]string{"msg": "You need to provide a valid number as digits"})
		return
	}

	stmt, err := db.Prepare("INSERT INTO otp (id, name, key, time, digits) values(NULL, ?, ?, ?, ?)")
	if err != nil {
		r.JSON(http.StatusBadRequest, err.Error())
		return
	}
	res, err := stmt.Exec(otp.Name, otp.Key, otp.Time, otp.Digits)
	if err != nil {
		r.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
		return
	}
	id, err := res.LastInsertId()
	if err != nil {
		r.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
		return
	}

	otp.Id = int(id)
	totp := &gotp.TOTP{Secret: strings.ToUpper(otp.Key), IsBase32Secret: true, Length: uint8(otp.Digits), Period: uint8(otp.Time)}
	r.JSON(http.StatusOK, Retval{otp.Id, otp.Name, otp.Time, totp.Get()})
}

func deleteOtp(r render.Render, params martini.Params, db *sql.DB) {
	id := params["id"]

	stmt, err := db.Prepare("DELETE FROM otp WHERE id = ?")
	if err != nil {
		r.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
		return
	}

	res, err := stmt.Exec(id)
	if err != nil {
		r.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
		return
	}

	affect, err := res.RowsAffected()
	if err != nil {
		r.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
		return
	}

	if affect > 0 {
		r.JSON(http.StatusOK, map[string]string{"msg": "Otp successfully deleted"})
	} else {
		r.JSON(http.StatusNotFound, map[string]string{"msg": "No such Otp"})
	}
}
