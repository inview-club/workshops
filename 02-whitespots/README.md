# Воркшоп № 2. WhiteSpots и технологии DevSecOps

## Установка Auditor
1. Склонируйте репозиторий `https://gitlab.inview.team/whitespots-public-fork/auditor.git`
2. Выполните `docker compose up -d`
3. Перейдите по адресу `127.0.0.1:8080` и сгенерируйте новый Access Token, сохраните его.
4. Добавьте этот токен в переменную `ACCESS_TOKEN` в .env файл в директории `auditor`.
5. Перезапустите контейнеры
    ```
    docker compose down
    docker compose up -d
    ```

## Установка AppSec Portal
1. Склонируйте репозиторий `https://gitlab.inview.team/whitespots-public-fork/appsec-portal.git`
2. Выполните `./set_vars.sh`
3. Добавьте версию `echo IMAGE_VERSION=release_v25.11.3 >> .env`
4. Запустите портал `sh run.sh`
5. Создайте пароль суперпользователя `docker compose exec back python3 manage.py createsuperuser --username admin`
6. Перейдите по адресу `127.0.0.1:80` и введите лицензионный ключ, который вы можете получить здесь: ([@inview_bot](https://t.me/inview_bot))
7. Перейдите в раздел Auditor - Config. Укажите адрес аудитора `http://host.docker.internal:8080/` и ваш Access Token, полученный ранее.
8. Измените Internal Portal URL на `http://host.docker.internal/`, кликнув на той же странице на Workflow Settings.
9. Добавьте приватный ключ (для клонирования репозиториев по SSH, должен быть привязан к вашему аккаунту на GitLab/GitHub/etc, не должен содержать пароля).

## Добавление репозиториев и сканирование
1. Перейдите во вкладку Assets - Repositories.
2. Добавьте новый репозиторий (репозитории с уязвимостями для тестирования - https://gitlab.com/whitespots-public/vulnerable-apps), указав URL для клонирования по SSH, добавьте теги.
3. Выберите репозиторий и запустите аудит.
4. После аудита просмотрите найденные уязвимости в Findings. Подтвердите (Verified) некоторые из них.
5. Попробуйте добавить свой проект и выполнить сканирование.


## Интеграция с IDE
1. Перейдите в настройки портала, выберите пункт IDE Integration.
2. Установите расширение для вашей IDE, в настройках расширения необходимо указать URL портала и API токен, который можно получить в профиле пользователя.
3. Склонируйте ранее добавленный репозиторий на локальную машину и откройте его в IDE
4. Проследите отображение уязвимостей в исходном коде в IDE.

## Интеграция с Gitlab CI
Для авторизации на инстансе [Gitlab](https://gitlab.inview.team) перейдите в бота и получите данные для входа.

Вам будет доступна группа `whitespots_<USERNAME>`

1. Выполните форк репозиториев pipelines, security-images и vulnerable-python-app из группы whitespots в свою группу.
2. Перейдите в репозиторий `security-images`, внесите изменения в файле `.gitlab-ci.yml` в разделе include - project, указав путь до репозитория pipelines в формате `whitespots_<USERNAME>/pipelines`.
3. Зарегистрируйте раннер через Settings - CI/CD - Runners. На вашем устройстве должен быть установлен gitlab-runner. Введите команду регистрации раннера. Тип раннера - docker. Образ - `alpine:latest`. Запустите раннер.
4. Создайте пайплайн через Build - Pipelines - New.
5. Выберите образы - toolset, gitleaks, trufflehog3. Запустите их сборку.
6. После сборки эти образы должны отображаться в Deploy - Container Registry проекта.
7. В вашем репозитории `pipelines` внесите изменения в файл `common/variables.yml`. Измените `SEC_PATH_TO_IMAGES` на путь, соответствующий вашему репозиторию с security-images.
8. Аналогичную замену произведите в файле `pipelines.yml` в пункте `image`.
9. Сгенерируйте Deploy Token в вашем репозитории `security-images` через Settings - Repository - Deploy Tokens.
10. Выполните кодировку в base64
`echo -n "TOKEN_NAME:PASSWORD" | base64`
11. В конфиге `.gitlab-runner/config.toml`, который должен располагаться в домашней директории вашего пользователя. В разделе `[[runners]]` найдите свой раннер и в разделе `[runners.docker]` добавьте следующее:
    ```
    environment = ["DOCKER_AUTH_CONFIG={\"auths\":{\"gitlab.inview.team:5050\":{\"auth\":\"YOUR_BASE64\"}}}"]
    extra_hosts = ["host.docker.internal:host-gateway"]
    ```
    **Обязательно перезапустите раннер!**

12. В вашем репозитории `vulnerable-python-app` создайте переменную  окружения `SEC_PORTAL_KEY` через Settings - CI/CD - Variables. Укажите свойство Masked, флаг Protect variable не требуется, снимите его.
13. Добавьте раннер в репозиторий (можно использовать уже существующий, просто активируйте его в настройках соответствующего репозитория - Enable for this project).
14. Создайте пайплайн (аналогично п. 4). Проследите за его выполнением. Некоторые проверки не будут выполнены из-за отсутствия соответствующих образов.
15. При успешном выполнении тестирования на портале появится новый отчёт (раздел Audits), а также репозиторий, если он ранее не был добавлен.

Если вы хотите запускать проверки репозиториев с нашего инстанса с использованием портала, вам необходимо добавить SSH ключ, а также в настройках пайплайна Code Downloader (В Аудиторе!) добавить следующее `ssh-keyscan -H gitlab.inview.team >> ~/.ssh/known_hosts`.

