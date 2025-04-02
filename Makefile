CC = gcc
SRC = lab-server-side/zombie.c
OUT = ./zombie

all:
	$(CC) $(SRC) -o $(OUT)
