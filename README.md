# mini_pos
Aplikasi mini_pos untuk test backend.

## Environment variable
Aplikasi membutuhkan variable yang wajib diisi untuk koneksi ke database dan keamanan JWT
.  
copy atau rename file `.env-example` menjadi `.env`. Sesuaikan isinya seperti pada contoh.
file ini akan diload ketika program dijalankan.

## Database
Aplikasi memerlukan database `PostgreSQL` dengan nama database `minipos`.  
adapun Table yang dibutuhkan ada pada file `doc/database.sql` : (tidak sempat dibuat automigration) berikut dengan ERD nya pada file `doc/minipos_erd.pdf`.

## Menjalankan Aplikasi
1. jalankan perintah `go mod tidy` untuk mendownload dependency
2. jalankan `go run main.go` untuk mulai menjalankan aplikasi semasa pengujian
3. buka browser dan jelajah `http://127.0.0.1:3500/swagger/index.html` untuk menjalankan dokumentasi rest-api
4. atau import file hasil export postman di folder /doc

## Swagger
untuk memperbarui doc swagger bisa menggunakan `swag init -g app/app.go` . Juga lakukan ini apabila isi folder /docs kosong.  
Swagger Doc bisa diakses melalui `http://127.0.0.1:3500/swagger/index.html`.

Gunakan token `Bearer<spasi><Token tanpa petik>` pada menu Authorize yang bisa dibuka dengan mengklik tombol hijau di kanan atas menu swager.


## Endpoint
postman config disertakan pada file `doc/minipos.postman_collection.json`.


### Daftar lengkap map url
```
	/*
	app.Static("/image/products", "./static/image/products")

	// url mapping
	api := app.Group("/api/v1")

	// Merchant Endpoint     << ---- pada merchant tidak memerlukan token untuk memudahkan pengetesan
	api.Post("/merchant", merchantHandler.CreateMerchant)
	api.Get("/merchant/:id", merchantHandler.GetMerchant)
	api.Get("/merchant", merchantHandler.FindMerchant)
	api.Put("/merchant/:id", merchantHandler.EditMerchant)
	api.Delete("/merchant/:id", merchantHandler.DeleteMerchant)

	// USER Endpont
	api.Get("/users/:id", userHandler.Get)
	api.Get("/users", userHandler.Find)
	api.Post("/login", userHandler.Login)
	api.Post("/refresh", userHandler.RefreshToken)
	api.Get("/profile", middleware.NormalAuth(), userHandler.GetProfile)
	api.Post("/register", middleware.FreshAuth(roles.RoleOwner), userHandler.Register)
	api.Put("/users/:id", middleware.NormalAuth(roles.RoleOwner), userHandler.Edit)
	api.Delete("/users/:id", middleware.NormalAuth(roles.RoleOwner), userHandler.Delete)

	// Outlet Endpont
	api.Get("/outlets/:id", middleware.NormalAuth(), outletHandler.Get)
	api.Get("/outlets", middleware.NormalAuth(), outletHandler.Find)
	api.Get("/current-outlet", middleware.NormalAuth(), outletHandler.GetCurrentOutlet)
	api.Post("/outlets", middleware.NormalAuth(roles.RoleOwner), outletHandler.CreateOutlet)
	api.Put("/outlets/:id", middleware.NormalAuth(roles.RoleOwner), outletHandler.Edit)
	api.Delete("/outlets/:id", middleware.NormalAuth(roles.RoleOwner), outletHandler.Delete)

	// Product Endpont
	api.Get("/products/:id", middleware.NormalAuth(), productHandler.Get)
	api.Get("/products", middleware.NormalAuth(), productHandler.Find)
	api.Post("/products", middleware.NormalAuth(roles.RoleOwner), productHandler.CreateProduct)
	api.Put("/products/:id", middleware.NormalAuth(roles.RoleOwner), productHandler.Edit)
	api.Delete("/products/:id", middleware.NormalAuth(roles.RoleOwner), productHandler.Delete)
	api.Post("/set-price", middleware.NormalAuth(roles.RoleOwner), productHandler.SetCustomPrice)
	api.Post("/products-image/:id", middleware.NormalAuth(roles.RoleOwner), productHandler.UploadImage)
	*/
```


## Memulai pengujian  <========================
1. Dimulai dari Merhcant endpoint, Pembuatan merchant akan membuat otomatis 1 user dengan role owner.  adapun passwordnya kita yang menentukan karena tidak ada verifikasi email pada aplikasi ini. gunakan email masukan dan password masukan sebagai data untuk login.
2. Ketika mulai login, user akan mendapatkan token JWT yang harus dibawa pada header dengan format Bearer. semua endpoint yang memiliki `middleware.NormalAuth()` akan mengecek keabsahan token dan role yang diperlukan. `middleware.FreshAuth()` memerlukan Token yang fresh (bukan hasil refresh token)
3. Buatlah satu buah outlet, outlet tersebut ditandai sebagai milik merchant yang sesuai dengan akun dengan role owner yang login.
4. Product memiliki data master harga yang agak unik perlakuannya. Menambahkan produk akan menambahkan master produk sesuai merhcant user.
5. User dapat menambahkan custom harga produk untuk outlet tertentu. untuk mendapatkan harga sesuai outlet tertentu, ketika melakukan get product harus menyertakan query `<url>?outlet=nomor_outlet`. contoh `{{url}}/api/v1/products/6?outlet=2`.  begitu juga dengan mendapatkan list product `{{url}}/api/v1/products?search=&outlet=2`. tanpa query outlet maka data master harga yang akan ditampilkan.


## Kontrak Struktur

### Middleware > Handler > Service > Dao || Api

- Handler digunakan untuk mengekstrak inputan dari user. params, query, json body, claims dari jwt serta validasi input
  ,termasuk memastikan dan menimpa huruf besar atau kecil.
- Service digunakan untuk bisnis logic, menggabungkan dua atau lebih dao atau utilitas pembantu lainnya, mengisi data
  yang dibutuhkan dao misalnya saat perpindahan dari requestData ke Model Data.
- Dao berkomunikasi langsung ke database. Beberapa kasus juga memastikan inputan huruf besar dan kecil pada inputan
  database yang caseSensitif untuk memaksimalkan indexing, memastikan nilai yang di input array<T> apabila array nil.
- Struktur tersebut memudahkan untuk pengetesan, terutama mock, namun pada tes ini membuat unit test akan memakan banyak waktu.
