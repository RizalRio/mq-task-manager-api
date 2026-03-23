-- Urutan DROP sangat penting! 
-- Hapus tabel yang memiliki Foreign Key (tasks) terlebih dahulu sebelum menghapus tabel induk (users)
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS users;