import unittest
from typing import List, Callable

reserved_characters = ['"', ':', '[', ',', '{', '}', ']']
ALL_OTHER_CHARS = '*'
SKIPPED_CHARS = [' ', '\n']

STACK: List[str] = []

START_TOKEN = False
TOKEN = ""


class StateMachineNode:
    def __init__(self, input_str: str, func: Callable[[str], None]=None):
        self.input_str_value = input_str
        self.transitions = {}
        self.func = func

    def add_transitions(self, input_str: str, node: 'StateMachineNode'):
        self.transitions[input_str] = node

    def is_valid_transition(self, input_str: str):
        if input_str not in reserved_characters and '*' in self.transitions:
            return True
        elif self.input_str_value == '*' and input_str not in reserved_characters:
            return True
        else:
            return input_str in self.transitions

    def _normalize_char(self, input_str: str):
        if input_str not in reserved_characters:
            return '*'
        return input_str

    def transition(self, input_str: str):
        if input_str in SKIPPED_CHARS:
            return self
        if not self.is_valid_transition(input_str):
            return None
        else:
            char = self._normalize_char(input_str)
            new_transition: 'StateMachineNode' = self.transitions[char]
            if new_transition.func:
                try:
                    new_transition.func(input_str)
                except Exception as e:
                    return False
            return new_transition


BEGIN_STATE: StateMachineNode = None


def enable_token(input_str: str):
    global START_TOKEN
    START_TOKEN = True
    global TOKEN
    TOKEN = ""
    TOKEN += input_str


def disable_token(input_str: str):
    validate_token(input_str)
    global START_TOKEN
    START_TOKEN = False
    global TOKEN
    TOKEN = ""


def collect_token(input_str: str):
    if input_str == '"' and START_TOKEN:
        return
    global TOKEN
    TOKEN += input_str


def validate_token(input_str: str):
    if not START_TOKEN:
        return
    global TOKEN
    if TOKEN.startswith(":"):
        TOKEN = TOKEN[1:]
    if not TOKEN:
        return
    if TOKEN == 'true' or TOKEN == 'false' or TOKEN == 'null' or TOKEN.startswith('"'):
        return
    try:
        int(TOKEN)
    except Exception as e:
        raise e
    return


def insert_in_stack(input_str: str):
    STACK.append(input_str)


def verify_in_stack(input_str: str):
    if input_str == '}':
        if STACK[len(STACK)-1] != "{":
            raise "Invalid json"
    elif input_str == ']':
        if STACK[len(STACK)-1] != "[":
            raise "Invalid json"
    STACK.pop()


def setup_transitions():
    begin_state = StateMachineNode('')
    global BEGIN_STATE
    BEGIN_STATE = begin_state

    json_start_node = StateMachineNode("{", insert_in_stack)
    token_start_state = StateMachineNode('"', enable_token)
    token_end_state = StateMachineNode('"', disable_token)
    comma_state = StateMachineNode(',', disable_token)
    list_start_state = StateMachineNode('[', insert_in_stack)
    list_end_state = StateMachineNode(']', verify_in_stack)
    end_state = StateMachineNode('}', verify_in_stack)
    key_state = StateMachineNode('*', collect_token)
    key_value_separator = StateMachineNode(':', enable_token)

    begin_state.add_transitions('{', json_start_node)
    json_start_node.add_transitions('"', token_start_state)
    json_start_node.add_transitions("}", end_state)

    token_start_state.add_transitions('*', key_state)

    key_state.add_transitions('*', key_state)

    key_state.add_transitions('"', token_end_state)
    key_state.add_transitions(',', comma_state)
    key_state.add_transitions("}", end_state)

    token_end_state.add_transitions('}', end_state)
    end_state.add_transitions("]", list_end_state)
    end_state.add_transitions(",", comma_state)

    token_end_state.add_transitions(',', comma_state)
    token_end_state.add_transitions(']', list_end_state)
    list_end_state.add_transitions('}', end_state)
    list_end_state.add_transitions(',', comma_state)

    comma_state.add_transitions('"', token_start_state)
    comma_state.add_transitions('{', json_start_node)
    comma_state.add_transitions('[', list_start_state)

    list_start_state.add_transitions("{", json_start_node)
    list_start_state.add_transitions("]", list_end_state)
    list_start_state.add_transitions('"', token_start_state)

    token_end_state.add_transitions(':', key_value_separator)

    key_value_separator.add_transitions('[', list_start_state)
    key_value_separator.add_transitions('"', token_start_state)
    key_value_separator.add_transitions('{', json_start_node)
    key_value_separator.add_transitions('*', key_state)


def test_json(input_str: str) -> bool:
    current_state: StateMachineNode = BEGIN_STATE
    index: int = 0
    valid_json = True
    while index < len(input_str):
        c = input_str[index]
        new_state: StateMachineNode = current_state.transition(c)
        if not new_state:
            valid_json = False
            break
        index += 1
        current_state = new_state
    if valid_json:
        return True
    else:
        return False


def get_file_content(path: str) -> str:
    from pathlib import Path
    return Path(path).read_text()


class TestStringMethods(unittest.TestCase):

    def setUp(self):
        setup_transitions()

    def test_empty(self):
        self.assertTrue(test_json(''), False)

    def test_simple_json(self):
        self.assertTrue(test_json('{}'), False)

    def test_simple_json_one_key_value(self):
        self.assertTrue(test_json('{\"abc\": \"def\"}'), False)

    def test_simple_invalid_json(self):
        self.assertFalse(test_json('{\"abc\": \"def\"}}'), False)
        self.assertFalse(test_json('{{\"abc\": \"def\"}'), False)
        self.assertFalse(test_json('{\"abc\":: \"def\"}'), False)
        self.assertFalse(test_json('{\"abc\"\": \"def\"}'), False)
        self.assertFalse(test_json('{\"abc\"\" \"def\"}'), False)

    def test_simple_json_with_list(self):
        self.assertTrue(test_json('{\"abc\": [\"def\", \"ghi\"]}'), True)

    def test_step1_invalid(self):
        content = get_file_content('tests/step1/invalid.json')
        self.assertTrue(test_json(content), False)

    def test_step1_valid(self):
        content = get_file_content('tests/step1/valid.json')
        self.assertTrue(test_json(content), True)

    def test_step2_invalid(self):
        content = get_file_content('tests/step2/invalid.json')
        self.assertFalse(test_json(content))

    def test_step2_invalid2(self):
        content = get_file_content('tests/step2/invalid2.json')
        self.assertFalse(test_json(content))

    def test_step2_valid(self):
        content = get_file_content('tests/step2/valid.json')
        self.assertTrue(test_json(content))

    def test_step2_valid2(self):
        content = get_file_content('tests/step2/valid2.json')
        self.assertTrue(test_json(content))

    def test_step3_invalid(self):
        content = get_file_content('tests/step3/invalid.json')
        self.assertFalse(test_json(content))

    def test_step3_valid(self):
        content = get_file_content('tests/step3/valid.json')
        self.assertTrue(test_json(content))

    def test_step4_invalid(self):
        content = get_file_content('tests/step4/invalid.json')
        self.assertFalse(test_json(content))

    def test_step4_valid(self):
        content = get_file_content('tests/step4/valid.json')
        self.assertTrue(test_json(content))

    def test_step4_valid2(self):
        content = get_file_content('tests/step4/valid2.json')
        self.assertTrue(test_json(content))


if __name__ == "__main__":
    setup_transitions()
    input_str = "{ \"key\": \"value\", \"key-n\": 101, \"key-o\": {}, \"key-l\": []}"
    print(test_json(input_str))