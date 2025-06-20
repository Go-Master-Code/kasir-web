Deskripsi sistem:
Berikut adalah mesin kasir berbasis web menggunakan Go (Golang) sebagai backend, dan JavaScript pada frontend.
Aplikasi ini akan digunakan di minimarket untuk mempermudah pengelolaan barang, transaksi, dan stok.
Berikut adalah fitur-fitur yang tersedia:

1. **Menu Dashboard**: 
    - Menunjukkan grafik jumlah pengunjung bulanan

2. **Menu Inventory (CRUD)**: 
    - Tambah barang baru.
    - Lihat daftar barang beserta detailnya (kode barang, nama barang, harga, stok).
    - Edit data barang tertentu.
    - Hapus barang dari inventory.

3. **Menu Gudang**:
    - Cek stok barang per kategori.
    - Tambah stok ke barang tertentu.
    - Lihat barang dengan stok rendah.

4. **Menu Transaksi**:
    - Buat transaksi baru dengan memasukkan kode barang dan jumlah.
    - Hitung total harga secara otomatis.
    - Cetak struk transaksi sederhana.

5. **Menu Laporan**:
    - Tampilkan laporan penjualan harian.
    - Tampilkan total pendapatan dalam periode tertentu.
    - Simpan laporan ke file.

6. **Menu Import Data Barang**:
    - Fitur untuk import data barang dari file excel yang akan diinput secara otomatis menggunakan goroutine ke dalam database
 
7. **Autentikasi Pengguna**:
    - Terdapat login untuk admin dan kasir.
    - Admin memiliki akses penuh ke semua fitur.
    - Kasir hanya dapat mengakses fitur transaksi dan laporan penjualan.
