# Воркшоп № 3. Git

## 0. Вступление

- Что за клуб
- Чем занимаемся

Тема воркшопа - Git.

Git - это система контроля версий.

Он помогает не потерять изменения и понимать, что происходило с проектом в определенный период времени.

Git - это инструмент для коллаборации.

Он позволяет удобно работать над репозиторием вместе с другими людьми.

Git - распределённая система.

Это значит, что у нас локально лежит не просто текущая версия файлов, а полноценный репозиторий со всей историей.

---

## 1. Состояние директории - 15 мин

### Теория

Commit — это сохранённое состояние проекта в конкретный момент времени.

Важно: Git хранит не «список изменений в файле», а снимок состояния файлов.

Это не означает, что Git каждый раз полностью копирует весь проект. Если содержимое файла не изменилось, Git переиспользует уже существующий объект. Поэтому хранение устроено достаточно эффективно.

Сам commit хранит:

- ссылку на состояние файлов;
- ссылку на предыдущий commit;
- автора;
- дату;
- commit message.

todo: вставить пример коммита

Порядок сохранения изменений

```text
Working tree
    ↓ git add
Index
    ↓ git commit
Repository
```



Команды:

```bash
git init
git status
git add
git diff
git diff --staged
git commit
git commit --amend
```

Работа `.gitignore` и зачем он нужен.

> `amend` создаёт новый commit, а не редактирует старый.

### Практика

* создать `config.txt`;
* создать `debug.log`;
* добавить `config.txt` в commit;
* игнорировать `*.log`;
* сделать commit;
* изменить `config.txt`;
* добавить изменение в предыдущий commit через `--amend`.

---

## 2. Ветки и stash — 15 мин

### Теория

Ветка — динамический указатель на конкретный commit

Когда мы делаем новый commit в ветку, указатель автоматически передвигается на него.

```text
A---B---C
        ^
       main
```

После создания ветки feature:

```text
A --- B --- C
          ^
         main
          ^
       feature
```

После добавления коммита в feature:
```text
A --- B --- C
          ^    \
         main   D
                ^
             feature
```


Команды:

```bash
git switch
git switch -c
git branch
```

Иногда нужно срочно переключиться на другую задачу, но текущие изменения ещё рано коммитить. В таком случае можно воспользоваться командой `git stash`, чтобы временно спрятать незаконченные изменения.

Команды:

```bash
git stash
git stash pop
```

### Практика

Вы в `feature/profile`.

Есть:
- изменённый `README.md`;
- новый файл `avatar.txt`.

Нужно:
1. временно спрятать изменения;
2. перейти в `main`;
3. создать `hotfix.txt` и закоммитить его;
4. вернуться в `feature/profile`;
5. восстановить изменения.
6. Закоммитить.

### Демо

```bash
git stash -u
git switch main

# hotfix

git switch feature
git stash pop
```


---

## 3. Merge и rebase — 20 мин

### Теория

Есть две ветки:

```text
      C---D feature
     /
A---B---E main

И main, и feature ушли вперёд. Нужно объединить изменения.

Merge:

```text
      C---D
     /     \
A---B---E---M
```

Создаётся дополнительный Merge commit, сливающий изменения двух веток.

Rebase:
rebase переносит коммиты одной ветки на коммиты другой (будто бы изменения были сделаны в одной ветке).

Было:

```text
      C---D
     /
A---B---E
```

Стало:

```text
A---B---E---C'---D'
```

C' и D' — новые commits с новыми hash.

Конфликты

В `main` ветке изменился файл и стал содержать `color=blue`

В другой ветке `feature` `color=red`

Git попросит решить конфликт изменений и расставит маркеры:
```
<<<<<<< HEAD
color=blue
=======
color=red
>>>>>>> feature
```

Для решения конфликта нужно вручную изменить файл, оставив итоговый вариант, затем добавить файл с помощью `git add` и `git merge --continue`. Для отмены мерджа можно использовать `git merge --abort`. Или `git rebase --continue` `git rebase --abort`.

### Практика

В main файл greeting.txt содержит:

### Демо

```bash
git merge feature
```

При конфликте:

```bash
git status

# исправить файл

git add file.txt
git merge --continue
```

Отмена:

```bash
git merge --abort
```

Для rebase:

```bash
git rebase main
git rebase --continue
git rebase --abort
```

Тезис:

> Локальную историю rebase'ить нормально. Общую опубликованную — осторожно.

Ссылки:

* https://git-scm.com/docs/git-merge
* https://git-scm.com/docs/git-rebase

---

## 4. Откаты — 20 мин

### Теория

#### reset

```bash
git reset --soft HEAD~1
git reset HEAD~1
git reset --hard HEAD~1
```

Разница:

```text
              commit   index   files

--soft          +        -       -
--mixed         +        +       -
--hard          +        +       +
```

`+` — состояние откатывается.

#### revert

```bash
git revert <commit>
```

Создаёт новый commit, отменяющий старый.

Тезис:

> `reset` двигает указатель.
> `revert` сохраняет историю и создаёт новый commit.

#### reflog

```bash
git reflog
```

Показывает, куда раньше указывал `HEAD`.

### Практика

```bash
git reset --hard HEAD~2
```

Задача: вернуть потерянные commits.

### Демо

```bash
git reflog
git reset --hard HEAD@{1}
```

Тезис:

> Часто commit не удалён — потерян только указатель на него.

Ссылки:

* https://git-scm.com/docs/git-reset
* https://git-scm.com/docs/git-revert
* https://git-scm.com/docs/git-reflog

---

## 5. Интересные кейсы — 20 мин

### Interactive rebase

История:

```text
Add feature
fix
fix2
oops
```

```bash
git rebase -i HEAD~4
```

Показать:

```text
pick
reword
squash
fixup
drop
```

Практика: превратить несколько commits в один нормальный.

---

### Worktree

Две ветки одновременно без второго clone:

```bash
git worktree add ../hotfix -b hotfix
git worktree list
```

---

### Hooks

Git может запускать скрипты:

```text
pre-commit
commit-msg
pre-push
```

Мини-демо: запретить commit message `test1`.

---

### Submodule

Репозиторий может ссылаться на конкретный commit другого репозитория.

```bash
git submodule add <repo> vendor/lib
git submodule update --init --recursive
```

Не углубляться — показать кейс и проблему, которую решает.

Ссылки:

* https://git-scm.com/docs/git-rebase
* https://git-scm.com/docs/git-worktree
* https://git-scm.com/docs/githooks
* https://git-scm.com/docs/git-submodule

---

## Сквозная команда

```bash
git log --oneline --graph --decorate --all
```

Можно сделать alias:

```bash
git config --global alias.lg "log --oneline --graph --decorate --all"
```

И дальше использовать:

```bash
git lg
```

## Главные тезисы

* `git add` → working tree → index;
* `commit` → новый снимок;
* branch → указатель;
* `merge` → объединение историй;
* `rebase` → новые commits на другой базе;
* `reset` → движение указателя;
* `revert` → новый обратный commit;
* `reflog` → история движения указателей.
