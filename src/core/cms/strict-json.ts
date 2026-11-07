/** Bounded JSON parsing rejects duplicate members, prototype keys and deep structures before content reaches schemas. */
export function parseCmsJSON(text: string): unknown {
  const parsed: unknown = JSON.parse(text);
  const stack: (Set<string> | null)[] = [];
  for (let i = 0; i < text.length; i++) {
    const char = text[i];
    if (char === "{" || char === "[") {
      stack.push(char === "{" ? new Set<string>() : null);
      if (stack.length > 64) throw new Error("CMS_CONTENT_INVALID");
    } else if (char === "}" || char === "]") stack.pop();
    else if (char === '"') {
      const start = i++;
      for (; i < text.length; i++) {
        if (text[i] === "\\") i++;
        else if (text[i] === '"') break;
      }
      let next = i + 1;
      while (next < text.length && /\s/.test(text[next] ?? "")) next++;
      if (text[next] === ":") {
        const key: unknown = JSON.parse(text.slice(start, i + 1));
        const keys = stack.at(-1);
        if (
          typeof key !== "string" ||
          !keys ||
          keys.has(key) ||
          ["__proto__", "constructor", "prototype"].includes(key)
        )
          throw new Error("CMS_CONTENT_INVALID");
        keys.add(key);
      }
    }
  }
  return parsed;
}
