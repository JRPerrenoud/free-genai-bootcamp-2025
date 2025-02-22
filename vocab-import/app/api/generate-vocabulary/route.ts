import { NextResponse } from "next/server"
import { Groq } from "groq-sdk"
import { z } from "zod"

const vocabularySchema = z.array(
  z.object({
    english: z.string(),
    spanish: z.string(),
  }),
)

// Log the API key (be careful not to expose this in production)
console.log("GROQ_API_KEY:", process.env.GROQ_API_KEY ? "Set" : "Not set")

export async function POST(req: Request) {
  try {
    if (!process.env.GROQ_API_KEY) {
      console.error("GROQ_API_KEY is not set")
      return NextResponse.json({ error: "GROQ_API_KEY is not set" }, { status: 500 })
    }

    const { theme } = await req.json()

    if (!theme) {
      return NextResponse.json({ error: "Theme is required" }, { status: 400 })
    }

    const groq = new Groq({
      apiKey: process.env.GROQ_API_KEY,
    })

    const completion = await groq.chat.completions.create({
      messages: [
        {
          role: "system",
          content: "You are a helpful assistant that generates vocabulary word pairs in English and Spanish.",
        },
        {
          role: "user",
          content: `Generate a list of 10 vocabulary words and their translations based on the theme: "${theme}". Respond only with a JSON array where each object has "english" and "spanish" keys.`,
        },
      ],
      model: "mixtral-8x7b-32768",
      temperature: 0.5,
      max_tokens: 1000,
    })

    const response = completion.choices[0]?.message?.content
    if (!response) {
      throw new Error("No response from Groq")
    }

    try {
      const jsonResponse = JSON.parse(response)
      const validatedResponse = vocabularySchema.parse(jsonResponse)
      return NextResponse.json(validatedResponse)
    } catch (parseError) {
      console.error("Error parsing response:", parseError)
      return NextResponse.json(
        { error: "Failed to parse vocabulary response" },
        { status: 500 },
      )
    }
  } catch (error) {
    console.error("Error in generate-vocabulary:", error)
    return NextResponse.json(
      {
        error: "Failed to generate vocabulary",
        details: error.message,
      },
      { status: 500 },
    )
  }
}
