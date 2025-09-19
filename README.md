# Утилиты для работы с событиями геймпада

Набор утилит для работы с файлами событий геймпада, записанными с помощью `evemu-record`. Утилиты позволяют объединять и повторять последовательности нажатий для автоматизации игровых процессов.

## Установка

```bash
# Клонирование репозитория
git clone https://github.com/CrazyBunker/evemu-utills.git gamepad-events-tools
cd gamepad-events-tools

# Сборка утилит
go build -o evemu-merge ./cmd/evemu-merge
go build -o evemu-repeat ./cmd/evemu-repeat

# Или установка в систему
go install ./cmd/evemu-merge
go install ./cmd/evemu-repeat
```

## Использование

### 1. Запись событий геймпада

Сначала запишите события с помощью `evemu-record`:

```bash
# Найти устройство геймпада
evemu-describe

# Записать события (замените /dev/input/eventX на ваше устройство)
evemu-record /dev/input/eventX > my_events.txt
```

### 2. Повторение событий - `repeat_events`

```bash
# Базовое использование: 5 повторов, вывод в файл
repeat_events my_events.txt 5 repeated_events.txt

# Чтение из stdin, запись в stdout
cat my_events.txt | repeat_events - 3 -

# Только вход из файла, вывод в stdout
repeat_events my_events.txt 5 -

# Чтение из stdin, вывод в файл
cat my_events.txt | repeat_events - 3 output.txt
```

### 3. Слияние событий - `merge_events`

```bash
# Базовое использование: объединение двух файлов
merge_events events1.txt events2.txt merged_events.txt

# Чтение базового файла из stdin
cat base_events.txt | merge_events - add_events.txt output.txt

# Вывод в stdout
merge_events events1.txt events2.txt -

# Полный pipeline со stdin/stdout
cat base.txt | merge_events - additions.txt - | repeat_events - 2 final.txt
```

### 4. Воспроизведение событий

```bash
# Воспроизведение с помощью evemu-play
evemu-play /dev/input/eventX < final_events.txt

# Или с использованием перенаправления
cat final_events.txt | evemu-play /dev/input/eventX
```

## Полный рабочий процесс

1. **Запись исходных событий**:
   ```bash
   evemu-record /dev/input/event4 > jump.txt
   evemu-record /dev/input/event4 > attack.txt
   ```

2. **Обработка событий**:
   ```bash
   # Повтор прыжка 3 раза
   repeat_events jump.txt 3 jump_x3.txt
   
   # Объединение атаки и прыжка
   merge_events attack.txt jump_x3.txt combo.txt
   
   # Повтор комбо 2 раза
   repeat_events combo.txt 2 final_combo.txt
   ```

3. **Воспроизведение**:
   ```bash
   cat final_combo.txt | evemu-play /dev/input/event4
   ```

## Примеры файлов событий

Файл событий в формате evemu:
```
# EVEMU 1.3
# Input device name: "Microsoft X-Box 360 pad"
################################
E: 0.000001 0003 0011 -001
E: 0.000001 0000 0000 0000
E: 0.583966 0003 0011 0000
E: 0.583966 0000 0000 0000
E: 1.355965 0001 0131 0001
```

## Опции командной строки

### `repeat_events`
```
repeat_events [входной_файл] <количество_повторов> [выходной_файл]

  входной_файл     - путь к файлу или '-' для stdin
  количество_повторов - число повторений последовательности
  выходной_файл    - путь к файлу или '-' для stdout
```

### `merge_events`
```
merge_events [базовый_файл] <добавочный_файл> [итоговый_файл]

  базовый_файл    - путь к файлу или '-' для stdin
  добавочный_файл - путь к файлу с событиями для добавления
  итоговый_файл   - путь к файлу или '-' для stdout
```

### Проблема: "События не воспроизводятся"
```bash
# Убедитесь, что файл содержит события
head -n 10 your_events.txt
```

## Примечания

- Утилиты сохраняют временные интервалы между событиями
- При объединении файлов временные метки автоматически корректируются
- Файлы событий совместимы с стандартным форматом evemu

## Зависимости

- Go 1.21 или новее
- Утилиты evemu (evemu-record, evemu-play)
- Доступ к устройствам ввода (/dev/input/)

## Лицензия

Проект распространяется под лицензией MIT.