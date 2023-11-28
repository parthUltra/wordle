import os
import random
from argparse import ArgumentParser

from colorama import Fore

with open("./words.txt") as wordle_words_file:
    wordle_words = wordle_words_file.read().splitlines()

with open("./all_words.txt") as all_words_file:
    all_words = all_words_file.read().splitlines()


def random_word():
    return random.choice(wordle_words)


def check(word: str):
    word = word.upper()
    if len(word) != 5:
        return False
    if word not in all_words:
        return False
    else:
        return True


def hard_check(word: str, keyboard_state):
    for i in word:
        if keyboard_state[i] == Fore.RED:
            return False
    return True


def prnt(word, inp):
    green = []
    yellow = []
    red = []
    inp = inp.upper()
    mask = [Fore.WHITE] * 5
    for i in range(5):  # For green
        if inp[i] == word[i]:
            mask[i] = Fore.GREEN
            word = word[:i] + " " + word[i + 1 :]

    for i in range(5):  # For yellow
        for j in range(5):
            if inp[i] == word[j] and mask[i] != Fore.GREEN:
                mask[i] = Fore.YELLOW
                word = word[:j] + " " + word[j + 1 :]

    for i in range(len(mask)):
        if mask[i] == Fore.GREEN:
            green.append(inp[i])
        if mask[i] == Fore.YELLOW:
            yellow.append(inp[i])
        if mask[i] == Fore.WHITE:
            red.append(inp[i])

    [print(mask[i] + inp[i] + Fore.RESET, end="") for i in range(5)]

    return [green, yellow, red]


def play(hard_mode=False, max_attempts=6):
    wordle = random_word()
    words = []
    attempts = 0
    word = ""

    all_keys = "QWERTYUIOPASDFGHJKLZXCVBNM"
    keys = list(all_keys)
    vals = [Fore.WHITE] * 26
    keyboard_state = {keys: vals for (keys, vals) in zip(keys, vals)}
    edits = []

    os.system("clear")
    print(Fore.GREEN + "\t\tWORDLE\n\n" + Fore.RESET)

    while word.upper() != wordle and attempts != max_attempts:
        print("\t\t", end="")
        word = input("-----\r\t\t")
        while not (check(word)):
            print("\t\t", end="")
            word = input("Invalid input, try again:\n\n\t\t-----\r\t\t")
        if hard_mode:
            while not (hard_check(word.upper(), keyboard_state)):
                print("\t\t", end="")
                word = input(
                    "Impossible input (Hard Mode Active), try again:\n\n-----\r\t\t"
                )
        words.append(word)
        os.system("clear")
        print(Fore.GREEN + "\t\tWORDLE\n\n" + Fore.RESET)
        for i in range(len(words)):
            print("\t\t", end="")
            edits = prnt(wordle, words[i])
            print("\t\t\n\n", end="")
            # print(f"\t\t\t\t{max_attempts-i-1} Attempts remaining:\n\n")
        attempts += 1
        keyboard_state = print_keys(keyboard_state, edits)
        print("\n\n")

    if word.upper() != wordle:
        print(f"\n\n\t\tGAME OVER!\n\n\t  The word was: {wordle}\n\n")
    else:
        print("\n\n\t\tCongrats!!!\n\n")


def print_keys(keyboard_state, edits):
    for i in edits[0]:
        keyboard_state[i] = Fore.GREEN
    for i in edits[1]:
        if keyboard_state[i] != Fore.GREEN:
            keyboard_state[i] = Fore.YELLOW
    for i in edits[2]:
        if keyboard_state[i] != Fore.GREEN and keyboard_state[i] != Fore.YELLOW:
            keyboard_state[i] = Fore.RED

    space = "   "
    print(space, end="")
    for key, color in keyboard_state.items():
        print(color + key + "   " + Fore.RESET, end="")
        if key == "P":
            print("\n\n  " + space, end="")
        if key == "L":
            print("\n\n    " + space, end="")
    return keyboard_state


def get_parser() -> ArgumentParser:
    parser = ArgumentParser(prog="wordle", description="command-line wordle")
    parser.add_argument("--hard", action="store_true", help="play wordle in hard mode")
    parser.add_argument(
        "--attempts", "-a", type=int, help="number of attempts", default=6
    )

    return parser


def main():
    parser = get_parser()
    args = parser.parse_args()

    play(hard_mode=args.hard, max_attempts=args.attempts)


if __name__ == "__main__":
    main()
