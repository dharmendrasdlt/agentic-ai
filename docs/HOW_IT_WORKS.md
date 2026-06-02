# How the System Works - Plain English Guide

This explains what happens when you ask a question, in simple, human-friendly language.

---

## The Journey of Your Question

### 1. You Ask a Question

You type: **"How does the PDF explain pump maintenance?"**

This is just plain text. Nothing fancy yet.

---

### 2. Your Question Gets Broken Into Pieces (Tokenization)

The system breaks your question into small pieces called **tokens** (kind of like words, but more precise):

```
Original: "How does the PDF explain pump maintenance?"

Becomes: 
  - "How"
  - "does"
  - "the"
  - "PDF"
  - "explain"
  - "pump"
  - "maintenance"
  - "?"
```

Why? Because the AI model works with tokens, not full sentences. It needs to understand individual pieces.

---

### 3. Your Question Becomes a Number (Vectorization)

Here's the cool part: **Your question is turned into a list of numbers.**

The system uses an AI model (Ollama, specifically the embedding model) to convert your text into what's called a **vector** - basically a list of 768 numbers that represent the *meaning* of your question.

```
Your question:
"How does the PDF explain pump maintenance?"

Gets converted to:
[0.123, -0.456, 0.789, 0.234, -0.567, 0.891, ... (764 more numbers)]

These numbers represent "the question is about PDF content and pump maintenance"
in a mathematical way that the system can understand.
```

**Important:** This happens in a special space where similar ideas have similar numbers.

---

### 4. Search for Similar Ideas (Vector Search)

Now the system asks: **"In the vector database, what chunks are most similar to these 768 numbers?"**

It's like searching, but instead of matching words, it matches *meaning*.

The database looks through all the document chunks (pieces of your PDFs) and finds the 3 most similar ones:

```
Top 3 Matches Found:

1. Match Score: 0.998 (99.8% similar)
   Text: "Pump maintenance involves checking pressure gauge, 
          oil level, and seals. Steps: 1) Turn off power..."
   From: Page 45, Chapter 3

2. Match Score: 0.876 (87.6% similar)
   Text: "Regular maintenance reduces downtime by 40% and 
         extends equipment life..."
   From: Page 51, Chapter 4

3. Match Score: 0.755 (75.5% similar)
   Text: "Tools needed: wrench, gauge, sealant..."
   From: Page 67, Chapter 5
```

Notice the **Match Score**? That tells you how confident the system is that this is relevant.

---

### 5. Get the Actual Text (Chunk Retrieval)

The system now has the 3 most relevant pieces of text from your PDFs.

Each piece comes with metadata (extra information):
- **What it says** (the actual text)
- **Where it came from** (which document, page, chapter)
- **How relevant it is** (the match score)

---

### 6. Build a Smart Question (RAG - Retrieval Augmented Generation)

The system now creates a special prompt for the AI. It's like saying:

> "Here's a question, and here are the most relevant documents about that question.
> Use ONLY these documents to answer. If you can't find the answer here, say 'I don't know.'"

The complete prompt looks like:

```
SYSTEM INSTRUCTION:
"You are a helpful assistant. Answer the user's question using ONLY 
the context provided below. If the answer cannot be found in the 
context, state clearly that you do not know."

CONTEXT FROM DOCUMENTS:
--- Source 1 (Score: 99.8%) ---
"Pump maintenance involves checking pressure gauge, oil level, and seals.
Steps: 1) Turn off power, 2) Release pressure safely, 3) Inspect 
each component for wear or damage..."

--- Source 2 (Score: 87.6%) ---
"Regular maintenance reduces downtime by 40% and extends equipment life."

--- Source 3 (Score: 75.5%) ---
"Tools needed: wrench, gauge, sealant..."

USER QUESTION:
"How does the PDF explain pump maintenance?"
```

---

### 7. AI Generates the Answer (LLM)

The AI model (Gemma 4) reads this prompt and starts generating an answer.

It does this **one token at a time**, streaming the answer in real-time:

```
Stream of tokens arriving:

Token 1: "Based"
Token 2: " on"
Token 3: " the"
Token 4: " provided"
Token 5: " PDF"
Token 6: " excerpts"
...

Building into: "Based on the provided PDF excerpts..."
```

You see this happening live on your screen - the answer appears word by word, like typing.

---

### 8. You See the Answer + Sources

Once done, you get:

**The Answer:**
> "Based on the provided PDF excerpts, pump maintenance involves checking 
> the pressure gauge, verifying oil levels, and inspecting seals. The 
> process requires: 1) Turn off power, 2) Release pressure safely, 
> 3) Inspect each component for wear or damage. Regular maintenance can 
> reduce downtime by up to 40%."

**Plus the Sources:**
```
✓ Source 1: Page 45, Chapter 3 (99.8% match)
✓ Source 2: Page 51, Chapter 4 (87.6% match)  
✓ Source 3: Page 67, Chapter 5 (75.5% match)
```

You can click each source to see the full text it came from.

---

## What Makes This Smart

### 1. Vector Search is Different from Text Search
- Text search: "Does this chunk contain the words 'pump' and 'maintenance'?"
- Vector search: "Does this chunk have the *same meaning* as the question?"

This is why it finds relevant content even if the exact words don't match.

### 2. It Uses Your Documents First
Before even trying the internet, it searches your PDFs. This means:
- ✅ You get answers grounded in YOUR documents
- ✅ You get exact citations where the answer came from
- ✅ The AI can't make things up - it only uses what's in your PDFs

### 3. Streaming Gives You Feedback
Instead of waiting for the full answer, you see it appear word-by-word. This tells you:
- The system is working (not hung)
- It's thinking and generating (not stuck)
- The answer is coming

---

## When Something Goes Wrong

### "I do not know"
This means:
- Your documents don't contain information about this topic
- The system tried hard but couldn't find relevant matches
- The AI is being honest instead of making something up

### Low Match Scores
If the match scores are low (like 0.45), it means:
- The question is quite different from anything in the documents
- The answer might not be reliable
- The system tried its best but the documents might not cover this

---

## The Two Types of Models Working Together

### Embedding Model (Ollama gemma4:e2b)
- **Job:** Convert text into vectors (numbers)
- **When used:** 
  - At ingestion time: converts each PDF chunk into vectors
  - At query time: converts your question into vectors
- **Goal:** Put similar ideas in nearby spots in vector space

### Generation Model (Ollama gemma4:e4b)
- **Job:** Read context and generate human-readable answers
- **When used:** After search results are found
- **Goal:** Write a natural, helpful answer

---

## The Complete Journey (Visual)

```
You type: "How does the PDF explain pump maintenance?"
    ↓
System breaks it into tokens
    ↓
Converts to vector (768 numbers)
    ↓
Searches vector database
    ↓
Finds 3 most similar document chunks
    ↓
Takes those chunks + your question
    ↓
Asks AI: "Answer using only these sources"
    ↓
AI generates answer word-by-word
    ↓
You see: "Based on the provided PDF excerpts..." 
         (appears in real-time, plus sources)
```

---

## Key Takeaways

1. **Your question becomes numbers** - vectors help find meaning, not just words
2. **Your documents are searched first** - answers are grounded in YOUR content
3. **You get sources** - you can verify where the answer came from
4. **It streams in real-time** - you see answers appearing word-by-word
5. **Two different AIs work together** - one for searching, one for writing
6. **It's honest about uncertainty** - says "I don't know" if documents don't have the answer

---

## That's It!

That's the complete journey from "I have a question" to "here's your answer with sources."

Pretty cool, right? 🚀
