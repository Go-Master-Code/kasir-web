package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"
	"text/template"
	"time"

	"github.com/Go-Master-Code/kasir-web/config"
	"github.com/Go-Master-Code/kasir-web/models"

	"github.com/gorilla/sessions"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

var db = config.OpenConnectionMaster()

// List of handler function

func BootstrapTemplate(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(path.Join("views", "template.html"))

	if err != nil {
		log.Println(err.Error())                                           //error spesifik untuk developer
		http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
		return
	}

	err = tmpl.Execute(w, nil) //tampilkan file html (index.html) yang di parse di atas
	if err != nil {
		log.Println(err.Error())                                           //error spesifik untuk developer
		http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
		return
	}
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello handler"))
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	// Halaman HTML Not Found
	tmpl, err := template.ParseFiles(path.Join("views", "not_found.html")) //buka file not_found.html
	if err != nil {
		log.Println(err)
		return
	}

	err = tmpl.Execute(w, nil) //tampilkan file html (not_found.html) yang di parse di atas
	if err != nil {
		log.Println(err.Error())                                           //error spesifik untuk developer
		http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
		return
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) { //parameter handler func wajib seperti ini
	//body function
	log.Println(r.URL.Path) //jika ingin tahu url path nya lari kemana

	if r.URL.Path != "/" { //Jika route nya tidak empty string (/), maka tampilkan http not found
		http.NotFound(w, r)
		return //jika masuk ke kondisi ini, stop, jangan lanjut ke coding bawah
	}

	//w.Write([]byte("Welcome to home!")) //parameter berupa array of byte, masukan saja string di dalamnya

	//menampilkan halaman html melalui handler home
	//yang dipakai adalam html template, bukan text template
	//panggil 2 view: index.html dan layout.html karena keduanya sudah terhubung dengan file layout.html
	tmpl, err := template.ParseFiles(path.Join("views", "index.html"), path.Join("views", "layout.html"))

	if err != nil {
		log.Println(err.Error())                                           //error spesifik untuk developer
		http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
		return
	}

	//buat data statis berbentuk map
	data := map[string]interface{}{ //key: string, value: interface kosong (bisa tipe data string, int, atau apa saja)
		"title":   "POS System",
		"content": "POS System With GO-Lang and PHP",
		"score":   100,
	}
	//key pada data map di atas sangat penting karena menjadi acuan untuk ditampilkan datanya pada web

	//data map di atas dimasukkan dalam parameter kedua func tmpl.Execute
	err = tmpl.Execute(w, data) //tampilkan file html (index.html) yang di parse di atas
	if err != nil {
		log.Println(err.Error())                                           //error spesifik untuk developer
		http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
		return
	}
}

func BarangHandler(w http.ResponseWriter, r *http.Request) {
	//cetak username yang didapat dari session setelah validasiAJAX
	session, _ := store.Get(r, "session-name")

	// Ambil username dari session
	username, _ := session.Values["username"].(string)

	tmpl, err := template.ParseFiles(path.Join("views", "barang.html"), path.Join("views", "template.html")) //buka file product.html
	if err != nil {
		log.Println(err)
		return
	}

	barang := models.TampilkanBarang(db)           //ambil data barang
	kategori := models.TampilkanKategoriBarang(db) //ambil data kategori barang

	data := struct { //buat struct yang menampung data barang dan kategori barang
		Barangs    []models.Barang
		Kategories []models.KategoriBarang
		Username   string
	}{
		Barangs:    barang,
		Kategories: kategori,
		Username:   username,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err.Error())                                           //error spesifik untuk developer
		http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
		return
	}
}

func BarangTipisHandler(w http.ResponseWriter, r *http.Request) {
	//cetak username yang didapat dari session setelah validasiAJAX
	session, _ := store.Get(r, "session-name")

	// Ambil username dari session
	username, _ := session.Values["username"].(string)

	tmpl, err := template.ParseFiles(path.Join("views", "barang.html"), path.Join("views", "template.html")) //buka file product.html
	if err != nil {
		log.Println(err)
		return
	}

	barang := models.TampilkanBarangSedikit(db)    //ambil data barang (return value dari model barang berupa struct)
	kategori := models.TampilkanKategoriBarang(db) //ambil data kategori barang

	data := struct { //buat struct yang menampung data barang dan kategori barang
		Barangs    []models.Barang
		Kategories []models.KategoriBarang
		Username   string
	}{
		Barangs:    barang,
		Kategories: kategori,
		Username:   username,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err.Error())                                           //error spesifik untuk developer
		http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
		return
	}
}

// contoh
func KategoriBarang(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(path.Join("views", "kategori_barang.html"), path.Join("views", "template.html")) //buka file product.html
	if err != nil {
		log.Println(err)
		return
	}

	user := models.TampilkanUser(db)
	kategori := models.TampilkanKategoriBarang(db)
	idKategoriTerpilih := 3

	data := struct {
		Users            []models.User
		Kategories       []models.KategoriBarang
		KategoriTerpilih int
	}{
		Users:            user,
		Kategories:       kategori,
		KategoriTerpilih: idKategoriTerpilih,
	}

	log.Println("Kategori terpilih: " + strconv.Itoa(data.KategoriTerpilih))

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err.Error())                                           //error spesifik untuk developer
		http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
		return
	}
}

func Form(w http.ResponseWriter, r *http.Request) { //parameter handler func wajib seperti ini
	//handler Form ini hanya bisa menerima request berupa get
	if r.Method == "GET" {
		//jika methodnya get
		tmpl, err := template.ParseFiles(path.Join("views", "form.html"), path.Join("views", "layout.html")) //buka form untuk input barang
		if err != nil {
			log.Println(err)
			http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())                                           //error spesifik untuk developer
			http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}
		return
	}
	//selain request GET, akan error
	http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError)
}

func Process(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		//ambil data dari form, lakukan post sesuai tipe form : POST
		err := r.ParseForm()
		if err != nil {
			log.Println(err.Error())                                        //error spesifik untuk developer
			http.Error(w, "Ini bukan POST", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}

		barang := r.Form.Get("barang")     //ambil nilai dari input box dengan id=barang pada modal
		kategori := r.Form.Get("kategori") //ambil nilai dari input box pada modal
		harga := r.Form.Get("harga")

		hargaInt, _ := strconv.Atoi(harga)
		kategoriInt, _ := strconv.Atoi(kategori)

		models.TambahBarang(db, barang, hargaInt, kategoriInt)

		http.Redirect(w, r, "/barang", http.StatusMovedPermanently)
	}
}

func UpdateBarang(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		//log.Println("Masuk method update data barang")
		//ambil data dari form, lakukan post sesuai tipe form : POST
		err := r.ParseForm()
		if err != nil {
			log.Println(err.Error())                                        //error spesifik untuk developer
			http.Error(w, "Ini bukan POST", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}

		id := r.Form.Get("idBarang")
		//log.Println("ID barang: " + id)
		barang := r.Form.Get("barang") //ambil nilai dari input box dengan id=barang pada modal
		//log.Println("Nama barang: " + barang)
		kategori := r.Form.Get("kategori") //ambil nilai dari input box pada modal
		harga := r.Form.Get("harga")
		stok := r.Form.Get("stok")

		idInt, _ := strconv.Atoi(id)
		hargaInt, _ := strconv.Atoi(harga)
		stokInt, _ := strconv.Atoi(stok)
		kategoriInt, _ := strconv.Atoi(kategori)

		models.UpdateBarang(db, idInt, barang, hargaInt, stokInt, kategoriInt)

		http.Redirect(w, r, "/barang", http.StatusMovedPermanently)
	} else if r.Method == "GET" {
		//cetak username yang didapat dari session setelah validasiAJAX
		session, _ := store.Get(r, "session-name")

		// Ambil username dari session
		username, _ := session.Values["username"].(string)

		//menampilkan data barang yang mau diedit
		tmpl, _ := template.ParseFiles(path.Join("views", "update_barang.html"), path.Join("views", "template.html"))

		log.Println("Masuk method get untuk update barang")
		idBarang := r.URL.Query().Get("id")
		log.Println("ID barang dari URL: " + idBarang)

		//ambil data parmeter id dari URL
		barang, selectedIdCategory := models.TampilkanBarangById(db, idBarang)
		log.Println("ID kategori terpilih: " + strconv.Itoa(selectedIdCategory))
		kategori := models.TampilkanKategoriBarang(db) //ambil data kategori barang

		data := struct { //buat struct yang menampung data barang dan kategori barang
			Barangs            []models.Barang
			Kategories         []models.KategoriBarang
			SelectedCategoryID int
			Username           string
		}{
			Barangs:            barang,
			Kategories:         kategori,
			SelectedCategoryID: selectedIdCategory,
			Username:           username,
		}

		log.Println("Kategori diselect: " + strconv.Itoa(data.SelectedCategoryID))

		err := r.ParseForm()
		if err != nil {
			log.Println(err.Error())                                       //error spesifik untuk developer
			http.Error(w, "Ini bukan GET", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}

		//log.Println("ID barang: " + barang[0].ID)
		//log.Println(barang)
		err = tmpl.Execute(w, data)
		if err != nil {
			log.Println(err.Error())                                           //error spesifik untuk developer
			http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}
	}
}

func DeleteBarang(w http.ResponseWriter, r *http.Request) {
	idBarang := r.URL.Query().Get("id")
	log.Println("Masuk ke method delete")
	log.Println("ID delete barang dari URL: " + idBarang)
	//ambil data parmeter id dari URL

	models.DeleteBarang(db, idBarang)
	http.Redirect(w, r, "/barang", http.StatusMovedPermanently)
}

func TabelTransaksi(w http.ResponseWriter, r *http.Request) {
	//cetak username yang didapat dari session setelah validasiAJAX
	session, _ := store.Get(r, "session-name")

	// Ambil username dari session
	username, _ := session.Values["username"].(string)

	tmpl, err := template.ParseFiles(path.Join("views", "transaksi.html"), path.Join("views", "template.html")) //buka file table.html
	if err != nil {
		log.Println(err)
		return
	}

	barang := models.TampilkanBarangOrderByNama(db)

	data := struct { //buat struct yang menampung data barang dan kategori barang
		Barangs  []models.Barang
		Username string
	}{
		Barangs:  barang,
		Username: username,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err.Error())                                           //error spesifik untuk developer
		http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
		return
	}
}

func SaveTransaksi(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		//ambil data dari form, lakukan post sesuai tipe form : POST
		err := r.ParseForm()
		if err != nil {
			log.Println(err.Error())                                        //error spesifik untuk developer
			http.Error(w, "Ini bukan POST", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}

		//cetak username yang didapat dari session setelah validasiAJAX
		session, _ := store.Get(r, "session-name")

		// Ambil username dari session
		username, _ := session.Values["username"].(string)

		//Save Master Transaksi
		idTransaksi := models.SaveMasterTransaksi(db, username)
		log.Println("ID Transaksi: " + idTransaksi + " telah tersimpan!")

		var kodeRow string = ""
		var qtyRow string = ""

		log.Println("Jumlah baris :" + r.Form.Get("jmlBaris"))
		//Ambil value jumlah iterasi dari form web (text box dengan id=jml)
		iterasi, _ := strconv.Atoi(r.Form.Get("jmlBaris"))

		//Iterasi sejumlah row barang pada web
		for i := 1; i <= iterasi; i++ {
			//i mulai dari 1 (settingan dari hlm web)
			kodeRow = r.Form.Get("kodeRow" + strconv.Itoa(i)) //untuk ambil parameter user dari web
			qtyRow = r.Form.Get("qtyRow" + strconv.Itoa(i))   //untuk ambil parameter user dari web

			if kodeRow != "" { //jika kode Row tidak kosong (ada valuenya)
				//parsing data ke int
				log.Println("Iterasi ke :" + strconv.Itoa(i))

				//Cetak kode dan qty barang per row
				log.Println("Kode barang :" + kodeRow)
				log.Println("Qty barang : " + qtyRow)

				kodeRowInt, _ := strconv.Atoi(kodeRow)
				qtyRowInt, _ := strconv.Atoi(qtyRow)

				//Masukkan data barang ke dalam struct (untuk jumlah row barang >1)
				models.TambahDetilTransaksi(idTransaksi, kodeRowInt, qtyRowInt)
			}
		}

		models.SaveDetilTransaksi(db) //Save detil transaksi dari data struct yang sudah terkumpul dari iterasi di atas

		//update stok barang setiap row detil transaksi
		//models.UpdateStokBarangDetilTransaksi(db)

		http.Redirect(w, r, "/transaksi", http.StatusMovedPermanently)
	}
}

func Container(w http.ResponseWriter, r *http.Request) { //parameter handler func wajib seperti ini
	//handler Form ini hanya bisa menerima request berupa get
	if r.Method == "GET" {
		//jika methodnya get
		tmpl, err := template.ParseFiles(path.Join("views", "container.html"), path.Join("views", "template.html")) //buka form untuk input barang
		if err != nil {
			log.Println(err)
			http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())                                           //error spesifik untuk developer
			http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}
		return
	}
	//selain request GET, akan error
	http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError)
}

func SearchDinamis(w http.ResponseWriter, r *http.Request) { //parameter handler func wajib seperti ini
	//handler Form ini hanya bisa menerima request berupa get
	if r.Method == "GET" {
		//jika methodnya get
		tmpl, err := template.ParseFiles(path.Join("views", "search_dinamis.html"), path.Join("views", "template.html")) //buka form untuk input barang
		if err != nil {
			log.Println(err)
			http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())                                           //error spesifik untuk developer
			http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}
		return
	}
	//selain request GET, akan error
	http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError)
}

func LaporanBarang(w http.ResponseWriter, r *http.Request) { //parameter handler func wajib seperti ini
	//handler Form ini hanya bisa menerima request berupa get
	if r.Method == "GET" {
		//jika methodnya get
		models.CetakLaporanStokBarang(db)

		tmpl, err := template.ParseFiles(path.Join("views", "laporan_barang.html"), path.Join("views", "template.html")) //buka form untuk input barang
		if err != nil {
			log.Println(err)
			http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())                                           //error spesifik untuk developer
			http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}
		return
	}
	//selain request GET, akan error
	http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError)
}

func LaporanTransaksiSummary(w http.ResponseWriter, r *http.Request) { //parameter handler func wajib seperti ini
	//handler Form ini hanya bisa menerima request berupa get
	if r.Method == "GET" {
		//jika methodnya get
		models.CetakLaporanTransaksiSummary(db)

		tmpl, err := template.ParseFiles(path.Join("views", "laporan_transaksi_summary.html"), path.Join("views", "template.html")) //buka form untuk input barang
		if err != nil {
			log.Println(err)
			http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())                                           //error spesifik untuk developer
			http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}
		return
	}
	//selain request GET, akan error
	http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError)
}

func LaporanTransaksiToday(w http.ResponseWriter, r *http.Request) { //parameter handler func wajib seperti ini
	//handler Form ini hanya bisa menerima request berupa get
	if r.Method == "GET" {
		//jika methodnya get
		models.CetakLaporanTransaksiSummaryToday(db)

		tmpl, err := template.ParseFiles(path.Join("views", "laporan_transaksi_today.html"), path.Join("views", "template.html")) //buka form untuk input barang
		if err != nil {
			log.Println(err)
			http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())                                           //error spesifik untuk developer
			http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}
		return
	}
	//selain request GET, akan error
	http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError)
}

func TampilkanPeriodeTransaksi(w http.ResponseWriter, r *http.Request) { //parameter handler func wajib seperti ini
	//handler Form ini hanya bisa menerima request berupa get
	if r.Method == "GET" {
		//jika methodnya get
		tmpl, err := template.ParseFiles(path.Join("views", "periode_transaksi.html"), path.Join("views", "template.html")) //buka form untuk input barang
		if err != nil {
			log.Println(err)
			http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())                                           //error spesifik untuk developer
			http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}
		return
	}
	//selain request GET, akan error
	http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError)
}

func LaporanTransaksiPeriode(w http.ResponseWriter, r *http.Request) { //parameter handler func wajib seperti ini
	//handler Form ini hanya bisa menerima request berupa POST
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Println(err.Error())                                        //error spesifik untuk developer
			http.Error(w, "Ini bukan POST", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}

		//jika methodnya POST
		//log.Println("Masuk method POST")

		//ambil nilai parameter dari dan sampai dari date picker di html

		dari := r.Form.Get("dari")
		sampai := r.Form.Get("sampai")
		//log.Println("Dari : " + dari + " Sampai : " + sampai)

		dariFull := dari + " 00:00:00"
		sampaiFull := sampai + " 23:59:59"
		//log.Println("Dari : " + dariFull + " Sampai : " + sampaiFull)

		models.CetakLapTransaksiSummaryPeriode(db, dariFull, sampaiFull, dari, sampai)

		tmpl, err := template.ParseFiles(path.Join("views", "laporan_transaksi_periode.html"), path.Join("views", "template.html")) //buka form untuk input barang
		if err != nil {
			log.Println(err)
			http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())                                           //error spesifik untuk developer
			http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}
		return
	}
	//selain request GET, akan error
	http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError)
}

func ImportBarang(w http.ResponseWriter, r *http.Request) { //import barang dari excel
	//handler Form ini hanya bisa menerima request berupa get
	if r.Method == "GET" {
		//jika methodnya get
		tmpl, err := template.ParseFiles(path.Join("views", "import_barang.html"), path.Join("views", "template.html")) //buka form untuk input barang
		if err != nil {
			log.Println(err)
			http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())                                           //error spesifik untuk developer
			http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}
		return
	}
	//selain request GET, akan error
	http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError)
}

func SaveImportBarang(w http.ResponseWriter, r *http.Request) { //save data barang dari excel ke db
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form untuk menangani file upload
	err := r.ParseMultipartForm(10 << 20) // Maksimum 10MB
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Ambil file dari form yang diupload
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Baca konten file yang di-upload
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusInternalServerError)
		return
	}

	// Buka file Excel menggunakan excelize
	f, err := excelize.OpenReader(bytes.NewReader(fileBytes))
	if err != nil {
		http.Error(w, "Unable to open Excel file", http.StatusInternalServerError)
		return
	}

	// Ambil sheet pertama dari file Excel
	sheet := f.GetSheetName(0)
	rows, err := f.GetRows(sheet)
	if err != nil {
		http.Error(w, "Unable to read rows", http.StatusInternalServerError)
		return
	}

	var dataList []models.Barang

	// Proses setiap baris data
	for i, row := range rows {
		if i == 0 { // Lewati baris header
			continue
		}

		if len(row) >= 5 { // Pastikan ada setidaknya 5 kolom
			harga, err := strconv.Atoi(row[2]) // Kolom ke-3 (harga)
			if err != nil {
				log.Printf("Invalid harga at row %d: %s\n", i+1, row[2])
				continue
			}

			stok, err := strconv.Atoi(row[3]) // Kolom ke-4 (stok)
			if err != nil {
				log.Printf("Invalid stok at row %d: %s\n", i+2, row[3])
				continue
			}

			id_kategori, err := strconv.Atoi(row[4]) // Kolom ke-5 (id kategori)
			if err != nil {
				log.Printf("Invalid ID kategori at row %d: %s\n", i+3, row[4])
				continue
			}

			barangs := models.Barang{
				NamaBarang: row[1], // Kolom ke-2 (nama barang)
				Harga:      harga,  // ambil dari strconv.atoi di atas
				Stok:       stok,
				IdKategori: id_kategori, // Kolom ke-5 (id kategori),
			}

			dataList = append(dataList, barangs)
		}
	}

	var db = config.OpenConnectionMaster()

	//manjalankan goroutine
	for _, barang := range dataList {
		go InsertBarangGoroutine(db, barang)
	}

	//berikan waktu jeda sampai semua goroutine selesai dieksekusi
	time.Sleep(time.Second * 1)

	//setelah selesai insert, redirect ke halaman data barang
	http.Redirect(w, r, "/barang", http.StatusMovedPermanently)

	/*****Apabila hendak melihat hasil data yang diinsert**********

	//Deklarasi struct apabila hendak menyertakan lebih dari 1 jenis data / struct
	type PageData struct {
		Title  string //data tambahan di samping struct hasil import data barang dari excel
		Barang []models.Barang
	}

	// Data untuk template
	pageData := PageData{
		Title:  "Data Dari File Excel",
		Barang: dataList,
	}

	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Menulis hasil rendering template ke response
	tmpl.Execute(w, pageData)

	********************/

}

func InsertBarangGoroutine(db *gorm.DB, barang models.Barang) {
	err := db.Create(&barang).Error
	if err != nil {
		panic(err)
	}
	fmt.Printf("Barang %s inserted successfully\n", barang.NamaBarang)
}

type Data struct {
	Label string `json:"label"`
	Value int    `json:"value"`
}

func GrafikBarangHandler(w http.ResponseWriter, r *http.Request) { //function ini tidak perlu dijadikan end point, pada halaman dashboard.html sudah dieksekusi
	data := []Data{
		{"Januari", 30},
		{"Februari", 50},
		{"Maret", 70},
		{"April", 90},
		{"Mei", 110},
		{"Juni", 200},
		{"Juli", 150},
		{"Agustus", 180},
		{"September", 200},
		{"Oktober", 100},
		{"November", 80},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func Dashboard(w http.ResponseWriter, r *http.Request) { //parameter handler func wajib seperti ini
	//handler Form ini hanya bisa menerima request berupa get
	if r.Method == "GET" {
		//cetak username yang didapat dari session setelah validasiAJAX
		session, _ := store.Get(r, "session-name")

		// Ambil username dari session
		username, _ := session.Values["username"].(string)
		//log.Println(username)

		// Kirim data ke template
		data := map[string]string{
			"Username": username,
		}

		//jika methodnya get
		tmpl, err := template.ParseFiles(path.Join("views", "dashboard.html"), path.Join("views", "template.html")) //buka form untuk input barang
		if err != nil {
			log.Println(err)
			http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			log.Println(err.Error())                                           //error spesifik untuk developer
			http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}
		return
	}
	//selain request GET, akan error
	http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError)
}

func Login(w http.ResponseWriter, r *http.Request) { //parameter handler func wajib seperti ini
	//handler Form ini hanya bisa menerima request berupa get
	if r.Method == "GET" {
		//jika methodnya get
		tmpl, err := template.ParseFiles(path.Join("views", "login.html")) //buka form untuk input barang
		if err != nil {
			log.Println(err)
			http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())                                           //error spesifik untuk developer
			http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}
		return
	}
	//selain request GET, akan error
	http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError)
}

func LoginAJAX(w http.ResponseWriter, r *http.Request) { //parameter handler func wajib seperti ini
	//handler Form ini hanya bisa menerima request berupa get
	if r.Method == "GET" {
		//jika methodnya get
		tmpl, err := template.ParseFiles(path.Join("views", "ajax.html")) //buka form untuk input barang
		if err != nil {
			log.Println(err)
			http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())                                           //error spesifik untuk developer
			http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}
		return
	}
	//selain request GET, akan error
	http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError)
}

//========validasi login + session========

// deklarasi var store untuk session
var store = sessions.NewCookieStore([]byte("secret-key"))

func ValidasiAJAX(w http.ResponseWriter, r *http.Request) {
	type LoginResponse struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}

	if r.Method == http.MethodPost {
		// Mengambil data dari form login
		r.ParseForm()

		username := r.FormValue("username")
		password := r.FormValue("password")

		//log username n password
		//log.Println(username)
		//log.Println(password)

		// validasi input username dan password
		id, pwd := models.ValidasiUser(db, username, password)

		// Validasi login
		if username == id && password == pwd {
			// Jika login berhasil
			// Login berhasil, buat session
			session, _ := store.Get(r, "session-name")
			session.Values["authenticated"] = true
			session.Save(r, w)

			sessionUser, _ := store.Get(r, "session-name")
			session.Values["username"] = username
			sessionUser.Save(r, w)
			//log.Println(session) bernilai true

			//log.Println("Login sukses")
			response := LoginResponse{
				Success: true,
				Message: "Login berhasil!",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else {
			// Jika login gagal
			//log.Println("Login gagal")
			response := LoginResponse{
				Success: false,
				Message: "Username atau password salah.",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// Fungsi logout untuk menghapus sesi
func Logout(w http.ResponseWriter, r *http.Request) {
	// Menghapus session (logout)
	session, _ := store.Get(r, "session-name")
	session.Values["authenticated"] = false
	session.Save(r, w)

	//log.Println(session) bernilai false
	// Pindah ke halaman login
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func RequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")
		// Cek apakah sesi 'authenticated' ada
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			// Jika tidak, arahkan ke halaman login (handler login)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
