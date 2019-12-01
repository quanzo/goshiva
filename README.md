goshiva
=======

Executing a command queue in multiple threads.

Программа выполняет команды в несколько потоков. Команды выбираются из очереди.
Очередь - файл, в котором одна команда на одной строке.

По умолчаню, 5 одновременно запущенных команд и файл queue.txt в папке с
программой.

Во время исполнения очереди, ее пожно дополнять, дописывая команды в конец
файла-очереди. Удалять строки из очереди нельзя.

Parameters of CLI
-----------------

| param         | type              | sample                    | description                                                          |
|---------------|-------------------|---------------------------|----------------------------------------------------------------------|
| \-queue-file  | string            | \-queue-file=./queue.txt  | Имя файла с очередью команд                                          |
| \-max-process | integer           | \-max-process=5           | Максимальное количество одновременно запущенных команд               |
| \-wait-queue  | bool              | \-wait-queue=true         | Ждать или нет новых команд в файле после окончания команд в очереди. |
| \-sleep-time  | integer           | \-sleep-time=1000         | Задержка в миллисекунда после старта команды из очереди.             |
| \-messages    | bool              | \-messages=true           | Вывод сообщение программы.                                           |
| \-output      | string (filename) | \-output=./out.txt        | Файл для записи вывода команд из очереди.                            |
| \-cmd-output  | string (filename) | \-cmd-output=./cmdout.txt | Файл для записи сообщений программы.                                 |

 

Queue command
-------------

| cmd            | param    | description                           |
|----------------|----------|---------------------------------------|
| \#             |          | Comment lin                           |
| :quit          |          | End process queue                     |
| :exit          |          |                                       |
| :running-count |          | Echo count running command            |
| :queue         |          | Echo queue file name                  |
| :status        |          | Echo status program                   |
| :sleep         | 1000     | Sleep program. Time in msec           |
| :delay         | 1000     |                                       |
| :change-queue  | filename | Run another command queue             |
| :echo          |          | Print string in program output stream |
|                |          |                                       |

 

 
