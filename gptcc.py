import sys
from openai import OpenAI

client = OpenAI()

message = ' '.join(sys.argv[1:])

response = client.chat.completions.create(
  model="gpt-4",
  messages=[
    {"role": "system", "content": """
     Input is a commit message. Output is a commit message with a Conventional Commits (CC) prefix.
     
     If the input already has a CC prefix just return it, otherwise add it. Do not otherwise modify the input unless the prefix looks like a scope is already there.

     Do not output anything other than the possibly modified commit message.
     """},
    {"role": "user", "content": message},
  ]
)

replaced = response.choices[0].message.content

print(replaced)