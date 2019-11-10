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

 

 

Queue command
-------------

| cmd            | param    | description                 |
|----------------|----------|-----------------------------|
| \#             |          | Comment lin                 |
| :quit          |          | End process queue           |
| :exit          |          |                             |
| :running-count |          | Echo count running command  |
| :queue         |          | Echo queue file name        |
| :status        |          | Echo status program         |
| :sleep         | 1000     | Sleep program. Time in msec |
| :delay         | 1000     |                             |
| :change-queue  | filename | Run another command queue   |
|                |          |                             |
|                |          |                             |

 

 
