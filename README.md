# goGetSwitch

## Дата обновления: 30 декабря 2021

### Настройка
В качестве БД на данный момент планирую использовать **mongodb**.
<br>
Необходимые команды для первоначальной настройки mongodb (установка и включение "автозапуска" у mongodb-сервиса):
```
apt install mongodb
systemctl enable mongodb.service
systemctl status mongodb.service
```

Содержимое ```/etc/systemd/system/goGetSwitch.service```:
```
[Unit]
Description = "Bot"
After = network.target

[Service]
ExecStart = /root/goGetSwitch/executables/goGetSwitch
# Под root'ом запускать, конечно, не круто, но на данный момент
# на сервере ничего, в общем-то, нет
User=root
WorkingDirectory = /root/
TimeoutSec=40s

[Install]
WantedBy = multi-user.target
```

Содержимое ```/etc/systemd/system/goGetSwitch.timer```:
```
[Install]
WantedBy=default.target

[Unit]
Description=Run getW every minute

[Timer]
# Start minutely on all days except weekend
OnCalendar=Mon..Fri *-*-* *:*:00
AccuracySec=1us
```

Включение сервисов (systemctl):
```
systemctl enable goGetSwitch.service
systemctl enable goGetSwitch.timer
systemctl start goGetSwitch.timer
systemctl status goGetSwitch.service
systemctl status goGetSwitch.timer
```

Установка Go:
```
sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt -y install golang
go version
```

### Идея:
Программа, которая будет скачивать и анализировать данные с рынка.
Затем при необходимости (при соответствии выбранной автоматически комбинации)
делать ставку. И сохранять все результаты в БД.

А также планирую (не знаю, сделаю ли в этой же наборе пакетов/репозитории, но я думаю, что можно - просто 
оформить в виде доп. пакетов, мб даже доп. main) сделать поиск самой 
лучшей комбинации за последние N дней
<br> 
**и/или**
поиск самой лучшей комбинации каждую ночь (в период с 0 до 4, когда выбранный брокер не работает)
и впоследствии в течение дня ведение статистики по комбинации (обновление рез-тов в БД и так далее. 
Можно, кстати, завести отдельный параметр для комбинации, который будет показывать, какой % ставок существует в тот же
период, когда существует ещё 5 других ставок (то есть они как бы "перекрывают" друг друга. Это ограничене 
брокера - 5 ставок одновременно)).
