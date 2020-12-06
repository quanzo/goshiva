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

| param         | short param | type              | sample                    | пояснение                                                            | description                                                                                 |
|---------------|-------------|-------------------|---------------------------|----------------------------------------------------------------------|---------------------------------------------------------------------------------------------|
| \-queue-file  | \-q         | string            | \-queue-file=./queue.txt  | Имя файла с очередью команд                                          | Command queue file name                                                                     |
| \-max-process | \-p         | integer           | \-max-process=5           | Максимальное количество одновременно запущенных команд               | Maximum number of simultaneously running commands                                           |
| \-wait-queue  | \-w         | bool              | \-wait-queue=false        | Ждать или нет новых команд в файле после окончания команд в очереди. | Whether or not to wait for new commands in the file after the end of commands in the queue. |
| \-sleep-time  | \-s         | integer           | \-sleep-time=1000         | Задержка в миллисекунда после старта команды из очереди.             | Delay in milliseconds after starting a command from the queue.                              |
| \-messages    | \-m         | bool              | \-messages=true           | Вывод сообщение программы.                                           | Output the program message.                                                                 |
| \-output      | \-o         | string (filename) | \-output=./out.txt        | Файл для записи вывода команд из очереди.                            | File for recording command output from the queue.                                           |
| \-cmd-output  | \-co        | string (filename) | \-cmd-output=./cmdout.txt | Файл для записи сообщений программы.                                 | File for recording program messages.                                                        |

 

Queue command
-------------

Эти команды могут присутствовать в очереди.

| cmd            | param    | пояснение                                                                                      | description                                                                                            |
|----------------|----------|------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------|
| \#             |          | Комментарий                                                                                    | Comment line                                                                                           |
| :quit          |          | Завершить выполнение очереди.                                                                  | End process queue                                                                                      |
| :exit          |          |                                                                                                |                                                                                                        |
| :running-count |          | Количество выполняющихся команд.                                                               | Echo count running command                                                                             |
| :queue         |          | Вывести имя файла с очередью команд.                                                           | Echo queue file name                                                                                   |
| :status        |          | Вывести статус программы                                                                       | Echo status program                                                                                    |
| :sleep         | 1000     | Задержка в выполнении следующей команды. В миллисекундах.                                      | Sleep program. Time in msec                                                                            |
| :delay         | 1000     |                                                                                                |                                                                                                        |
| :change-queue  | filename | Программа переключается на другую очередь                                                      | Run another command queue                                                                              |
| :echo          |          | Просто печатает сообщение в потоке выводв программы                                            | Print string in program output stream                                                                  |
| :wait_complete |          | Прежде, чем пойти далее по очереди, программа ожидает исполнения всех ранее запущенных команд. | Before going further in turn, the program waits for the execution of all previously launched commands. |

 

Runtime command
---------------

The launched program accepts commands from the keyboard. All of the commands
above can be used while executing the program.

Запущенная программа принимает команды с клавиатуры. Все команды, указанные
выше, можно использовать при выполнении программы.

 
