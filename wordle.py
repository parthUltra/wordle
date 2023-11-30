#!/usr/bin/env python3
import os
import random
from argparse import ArgumentParser

from colorama import Fore

WORDLE_WORDS_FILE = "./words.txt"  # Put the location of the wordle_words_file here
ALL_WORDS_FILE = "./all_words.txt"  # Put the location of all legal words here


class WORDLE:
    def __init__(self):
        with open(WORDLE_WORDS_FILE) as wordle_words_file:
            self.wordle_word = random.choice(
                wordle_words_file.read().splitlines()
            )  # The wordle word is set

        with open(ALL_WORDS_FILE) as all_words_file:
            self.all_words = (
                all_words_file.read().splitlines()
            )  # The all words file is read into an array

    def check(self, word: str) -> bool:  # Checks is the word is valid
        word = word.upper()
        if len(word) != 5:
            return False
        if word not in self.all_words:
            return False
        else:
            return True

    def hard_check(
        self, word: str, keyboard_state
    ) -> bool:  # Hard mode conditions check
        for i in word:
            if keyboard_state[i] == Fore.RED:
                return False
        return True

    def prnt(
        self, word: str, inp: str
    ) -> list[list[str]]:  # Prints the word in the appropriate colors
        # and returns the changes to keyboard_state identified
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

    def play(self, hard_mode=False, max_attempts=6):  # Function to play wordle
        wordle = self.wordle_word
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
            while not (self.check(word)):
                print("\t\t", end="")
                word = input("Invalid input, try again:\n\n\t\t-----\r\t\t")
            if hard_mode:
                while not (self.hard_check(word.upper(), keyboard_state)):
                    print("\t\t", end="")
                    word = input(
                        "Impossible input (Hard Mode), try again:\n\n\t\t-----\r\t\t"
                    )
            words.append(word)
            os.system("clear")
            print(Fore.GREEN + "\t\tWORDLE\n\n" + Fore.RESET)
            for i in range(len(words)):
                print("\t\t", end="")
                edits = self.prnt(wordle, words[i])
                print("\t\t\n\n", end="")
                # print(f"\t\t\t\t{max_attempts-i-1} Attempts remaining:\n\n")
            attempts += 1
            keyboard_state = self.print_keys(keyboard_state, edits)
            print("\n\n")

        if word.upper() != wordle:
            print(f"\n\n\t\tGAME OVER!\n\n\t  The word was: {wordle}\n\n")
        else:
            print("\n\n\t\tCongrats!!!\n\n")

    def print_keys(
        self, keyboard_state: dict[str, str], edits: list[list[str]]
    ) -> dict[
        str, str
    ]:  # Updates the keyboard_state, prints the keyboard and returns the keyboard
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

    def get_parser(self) -> ArgumentParser:
        parser = ArgumentParser(prog="wordle", description="command-line wordle")
        parser.add_argument(
            "--hard", action="store_true", help="play wordle in hard mode"
        )
        parser.add_argument(
            "--attempts", "-a", type=int, help="number of attempts", default=6
        )

        return parser

    def main(self):
        parser = self.get_parser()
        args = parser.parse_args()

        self.play(hard_mode=args.hard, max_attempts=args.attempts)


if __name__ == "__main__":
    WORDLE().main()
