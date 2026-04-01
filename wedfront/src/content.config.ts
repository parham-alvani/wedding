import { defineCollection, z } from "astro:content";
import { glob } from "astro/loaders";

const showcase = defineCollection({
  loader: glob({ pattern: "**/*.json", base: "./src/content/showcase" }),
  schema: ({ image }) =>
    z.object({
      title: z.string().min(1),
      image: image(),
      url: z.string().url(),
      featured: z.number().min(1).optional(),
    }),
});

export const collections = {
  showcase,
};
