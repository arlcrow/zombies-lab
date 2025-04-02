#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/wait.h>

int main() {
    pid_t pid = fork();

    if (pid < 0) {
        // Ошибка при вызове fork
        perror("fork failed");
        exit(1);
    }

    if (pid == 0) {
        // Дочерний процесс
        printf("Child process %d exiting\n", getpid());
        exit(0); // Дочерний процесс завершается
    } else {
        // Родительский процесс
        printf("Parent process %d running, child is %d\n", getpid(), pid);
        sleep(3600); // Спим 1 час, оставляя дочерний процесс в состоянии зомби
    }

    return 0;
}
