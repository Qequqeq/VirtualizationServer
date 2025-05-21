# VirtualizationServer

Руководство по установке и настройке сервера виртуализации.

---

## 1. Требования к системе
- Alpine Linux 3.21 (рекомендуется)
- Доступ с правами root
- Базовые навыки работы с терминалом

---

## 2. Установка

### 2.1 Подготовка Alpine Linux
#### 2.1.1. Скачайте образ: [Официальный сайт](https://alpinelinux.org/downloads/). Подойдет standart или extended
#### 2.1.2. Войдите от root пользователя. Пароль не требуется, просто введите root и нажммите Enter
#### 2.1.3. Установите базовую систему:
   ```bash
   setup-alpine
```
Далее приведено руководство по базовой настройке системы.На экране Вы будете постепенно видеть опции для настройки, ниже приведены наиболее простые и естественные действия. Конечно, Вы можете настроить систему и каким-то другим образом.
Там, где ввод -- пустая строка, система выбирает вариант по умолчанию. Это соотвествует вводу Enter в списке ниже.
* us
* us
* Enter
* Enter
* Enter
* Enter
* Ваш пароль для рут пользователя (toor)
* Введите пароль снова
* Europe/Moscow
* Enter
* Enter
* Enter
* Enter
* Enter
* Enter
* Enter
* В выводе с доступными дисками выберите диск, на который хотите установить систему.
* sys
* y
Если Вы устанавливаете ОС на физический стенд, выполните команду
```bash
reboot
```
и после выключения машины достаньте установочную флешку.

Если Вы используете VirtualBox, используйте команду 
```bash
poweroff
```
для выключения машины. После этого перейдите в настройки ВМ и удалите диск-загрузчик, чтобы в системе остался только виртуальный жосткий диск. После этого запустите машину.
### 2.3 Подготовка системы к работе
Выполните команды для установки необходимых пакетов и загрузки репозитория.
```bash
apk update && apk add git vim bash
wget https://github.com/git-lfs/git-lfs/releases/download/v3.6.1/git-lfs-linux-amd64-v3.6.1.tar.gz
tar -xvf git-lfs-linux-amd64-v3.6.1.tar.gz
rm git-lfs-linux-amd64-v3.6.1.tar.gz
cd git-lfs-3.6.1/ && ./install.sh
cd
git clone https://github.com/Qequqeq/VirtualizationServer
cd VirtualizationServer/
git lfs pull
```

Если по каким то причинам lfs не работает, воспользуйтесь 
```bash
./if_git_not_work
```

Для того, чтобы вы могли скачать все необходимые пакеты, настройте пакетный менеджер apk.
```bash
vim /etc/apk/repositories
```
Уберите комментарий перед репозиторием, оканчивающимся на /community
Установите пакеты, которые используются в проекте командой
```bash
cat installed_packages.txt | xargs apk add
```
### 2.4 Настройка сервиса Tuna
Перейдите на [официальный сайт](https://tuna.am/) и зарегистрируйтесь на сайте, купите подписку и зайдите в свой профиль. Зарезирвируйте домен в меню Домены. Выполните команды, чтобы привязать устройство к Tuna
```bash
curl -sSLf https://get.tuna.am | sh
tuna config save-token <hereYourToken>
```
### 2.5 Финальная настройки
Выпоните команду 
```bash
./initialization.sh <yourDomain>
```
Для изменения домена на Ваш.

Скомпилируйте исходный код, воспользовавшись командой
```bash
cd vm-api/ && go build -o <filename>
```

Теперь настроим выполнение регулярных скриптов, используя планировщик cron. Выполните
```bash
crontab -e
```
Допишите строчки 
```bash
@reboot time sleep 5 && cd /root/VirtualizationServer/vm-api/scripts/ && ./start_tuna.sh
0 */4 * * * cd /root/VirtualizationServer/vm-api/scripts/ && ./check_qemu/sh
0 */4 * * * cd /root/VirtualizationServer/vm-api/scripts/ && ./clear_ports.sh
```
выйдите и перезапустите систему. После этого запустите сервис.
```bash
cd VirtualizationServer/vm-api
./<filename>
```

