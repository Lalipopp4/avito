Данный проект является тестовым заданием к отбору на стажировку в Авито.

Методы запуска проекта:
Файл запуска находится в cmd/app
1.  $ go build -o app ./cmd/app
    $ ./app
2.  Запуск Makefile

Используемые фреймворки:
1.  github.com/lib/pq
2.  github.com/redis/go-redis/v9
3.  github.com/gin-gonic/gin
Также использовались пакеты из стандартной библиотеки

Примеры запросов к сервису через Postman:
1.  ![Alt text](images/image.png)
2.  ![Alt text](images/image-1.png)
3.  ![Alt text](images/image-2.png)
4.  ![Alt text](images/image-3.png)
5.  ![Alt text](images/image-4.png)
6.  ![Alt text](images/image-5.png)
7.  ![Alt text](images/image-6.png)
8.  ![Alt text](images/image-7.png)
9.  ![Alt text](images/image-8.png)
10. ![Alt text](images/image-9.png)
11. ![Alt text](images/image-10.png)
12. ![Alt text](images/image-11.png)
13. ![Alt text](images/image-12.png)
14. ![Alt text](images/image-13.png)
15. ![Alt text](images/image-14.png)
16. ![Alt text](images/image-115.png)

База данных:
    users - список пользователей
    segments - список сегментов
    user_segments - список сегментов пользователей, где каждая запись - пользователь и сегмент, в котором он состоит
    segment_history - история добавления и удаления пользователей в сегмент

Дополнительные задания:
1.  Реализовано в методе /segmenthistory и /csvhistory. Проводится поиск по записям в таблице segment_history на              совпадение даты с введенной. Возвращается url, в котором хранится сохраненный файл в формате .csv. Запрос на данный адрес возвращает этот файл.   
2.  Реализована горутина (checkTTL), которая в отдельном потоке проверяет дату автоматического удаления из таблицы user_segments на совпадение с текущeй датой. При успехе удаляет эти строки, иначе только "засыпает" на сутки.
3.  Реализовано в методе добавления сегмента через набор случайных(насколько позволяет math/rand) пользователей
    и последующем их добавлении в сегмент.
