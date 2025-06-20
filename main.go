package main

import (
	"log"
	"net/http"

	"github.com/Go-Master-Code/kasir-web/handler"
)

//var store = sessions.NewCookieStore([]byte("secret-key"))

/*========Function experimental session========
func loginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//jika methodnya get
		tmpl, err := template.ParseFiles(path.Join("views", "login2.html")) //buka form untuk input barang
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

// Fungsi untuk memproses login
func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Validasi login (contoh sederhana, bisa diganti dengan database)
	if username == "username" && password == "password" {
		// Login berhasil, buat session
		session, _ := store.Get(r, "session-name")
		session.Values["authenticated"] = true
		session.Save(r, w)

		http.Redirect(w, r, "/dashboard2", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/login2", http.StatusSeeOther)
	}
}

// Fungsi untuk memastikan pengguna sudah login
func requireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")
		// Cek apakah sesi 'authenticated' ada
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			// Jika tidak, arahkan ke halaman login (root endpoint)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Fungsi untuk halaman dashboard (hanya bisa diakses setelah login)
func dashboard2(w http.ResponseWriter, r *http.Request) {
	//http.Redirect(w, r, "/barang", http.StatusSeeOther)
	fmt.Fprintf(w, "Selamat datang di dashboard, Anda telah login!")
}

// Fungsi logout untuk menghapus sesi
func logout(w http.ResponseWriter, r *http.Request) {
	// Menghapus session (logout)
	session, _ := store.Get(r, "session-name")
	session.Values["authenticated"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/login2", http.StatusSeeOther)
}
*/

func main() {
	mux := http.NewServeMux() //buat server mux
	//mux.HandleFunc("/", handler.HomeHandler) 									  //kalau URL yang diakses /, maka handler yang dipanggil adalah ini

	mux.HandleFunc("/", handler.NotFound)                                               //default URL akan mengakses halaman 404 Not Found
	mux.HandleFunc("/grafik", handler.GrafikBarangHandler)                              //mengirimkan data dalam bentuk json untuk frontend
	mux.Handle("/dashboard", handler.RequireLogin(http.HandlerFunc(handler.Dashboard))) //untuk mengakses dashboard harus lolos login dulu

	mux.HandleFunc("/login", handler.Login)                                                         //halaman login
	mux.HandleFunc("/logout", handler.Logout)                                                       //Proses logout
	mux.HandleFunc("/ajax", handler.ValidasiAJAX)                                                   //validasi login dengan AJAX
	mux.Handle("/barang", handler.RequireLogin(http.HandlerFunc(handler.BarangHandler)))            //handler untuk data barang (semua)
	mux.Handle("/barang_tipis", handler.RequireLogin(http.HandlerFunc(handler.BarangTipisHandler))) //handler untuk data barang (stok tipis)
	//mux.HandleFunc("/barang", handler.BarangHandler)
	//mux.HandleFunc("/barang_tipis", handler.BarangTipisHandler) //handler untuk data barang (semua)
	mux.HandleFunc("/kategori", handler.KategoriBarang)    //handler menampilkan kategori barang pada combobox
	mux.HandleFunc("/template", handler.BootstrapTemplate) //template utama

	mux.HandleFunc("/process", handler.Process)     //daftarkan method untuk insert barang
	mux.HandleFunc("/update", handler.UpdateBarang) //handler untuk update barang
	mux.HandleFunc("/delete", handler.DeleteBarang) //handler untuk delete barang
	mux.Handle("/transaksi", handler.RequireLogin(http.HandlerFunc(handler.TabelTransaksi)))
	//mux.HandleFunc("/transaksi", handler.TabelTransaksi)                          //handler untuk transaksi
	mux.HandleFunc("/save_transaksi", handler.SaveTransaksi) //handler untuk save master-detil transaksi
	mux.HandleFunc("/container", handler.Container)          //handler untuk contoh container / panel
	mux.Handle("/laporan_barang", handler.RequireLogin(http.HandlerFunc(handler.LaporanBarang)))
	//mux.HandleFunc("/laporan_barang", handler.LaporanBarang)                      //handler untuk search dinamis melalui text box (diberi suggestion berdasarkan input user)
	mux.Handle("/laporan_transaksi", handler.RequireLogin(http.HandlerFunc(handler.LaporanTransaksiSummary)))
	//mux.HandleFunc("/laporan_transaksi", handler.LaporanTransaksiSummary)         //handler untuk search dinamis melalui text box (diberi suggestion berdasarkan input user)
	mux.Handle("/laporan_transaksi_today", handler.RequireLogin(http.HandlerFunc(handler.LaporanTransaksiToday)))
	//mux.HandleFunc("/laporan_transaksi_today", handler.LaporanTransaksiToday)     //handler untuk generate dan menampilkan laporan penjualan hari ini
	mux.Handle("/periode_transaksi", handler.RequireLogin(http.HandlerFunc(handler.TampilkanPeriodeTransaksi)))
	//mux.HandleFunc("/periode_transaksi", handler.TampilkanPeriodeTransaksi)       //handler memasukkan periode laporan penjualan
	mux.Handle("/laporan_transaksi_periode", handler.RequireLogin(http.HandlerFunc(handler.LaporanTransaksiPeriode)))
	//mux.HandleFunc("/laporan_transaksi_periode", handler.LaporanTransaksiPeriode) //handler untuk cetak laporan penjualan per periode tertentu
	mux.Handle("/import_barang", handler.RequireLogin(http.HandlerFunc(handler.ImportBarang)))
	//mux.HandleFunc("/import_barang", handler.ImportBarang)                        //handler untuk import data barang dari excel
	mux.HandleFunc("/save_import_barang", handler.SaveImportBarang) //menyimpan data hasil import dari excel ke dalam database

	/*========Handler contoh dan experimen========
	mux.HandleFunc("/form", handler.Form)                                         //contoh form
	mux.HandleFunc("/hello", handler.HelloHandler)                                //contoh

	//mux.HandleFunc("/search", handler.SearchDinamis)                            //handler untuk search dinamis melalui text box (diberi suggestion berdasarkan input user)

	mux.HandleFunc("/login_ajax", handler.LoginAJAX)                              //tampil halaman login dengan AJAX
	mux.HandleFunc("/logout", handler.Logout)                                     // Proses logout
	mux.HandleFunc("/login2", loginPage) 										  // Tampilkan halaman login
	mux.HandleFunc("/auth", login)       										  // Proses login
	mux.HandleFunc("/logout2", logout)											  // Proses logout

	// Halaman yang hanya dapat diakses jika sudah login
	mux.Handle("/dashboard2", requireLogin(http.HandlerFunc(dashboard2)))
	*/

	//load file css dalam folder style
	fileServer := http.FileServer(http.Dir("assets")) //load directory
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	//buat log untuk tahu server sudah berjalan di port 3000
	log.Println("Server running on port 3000")

	//menjalankan server
	err := http.ListenAndServe("localhost:3000", mux)
	log.Fatal(err) //jika terjadi error
}
