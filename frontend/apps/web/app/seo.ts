const defaultOrigin = "https://referral.vivl.cc"

export const siteOrigin = (
  import.meta.env.VITE_SITE_URL || defaultOrigin
).replace(/\/$/, "")

export const absoluteURL = (path: string) =>
  `${siteOrigin}${path.startsWith("/") ? path : `/${path}`}`

export function seoMeta({
  title,
  description,
  path,
  index = true,
}: {
  title: string
  description: string
  path: string
  index?: boolean
}) {
  const url = absoluteURL(path)
  const image = absoluteURL("/og-image.png")
  return [
    { title },
    { name: "description", content: description },
    {
      name: "robots",
      content: index
        ? "index, follow, max-image-preview:large, max-snippet:-1, max-video-preview:-1"
        : "noindex, nofollow, noarchive",
    },
    { tagName: "link", rel: "canonical", href: url },
    { property: "og:type", content: "website" },
    { property: "og:url", content: url },
    { property: "og:title", content: title },
    { property: "og:description", content: description },
    { property: "og:image", content: image },
    { property: "og:image:type", content: "image/png" },
    { property: "og:image:width", content: "1200" },
    { property: "og:image:height", content: "630" },
    {
      property: "og:image:alt",
      content: "Referral invitation and rewards platform",
    },
    { name: "twitter:card", content: "summary_large_image" },
    { name: "twitter:title", content: title },
    { name: "twitter:description", content: description },
    { name: "twitter:image", content: image },
    {
      name: "twitter:image:alt",
      content: "Referral invitation and rewards platform",
    },
  ]
}

export const structuredData = {
  "@context": "https://schema.org",
  "@graph": [
    {
      "@type": "Organization",
      "@id": `${siteOrigin}/#organization`,
      name: "Referral",
      url: `${siteOrigin}/`,
      logo: absoluteURL("/apple-touch-icon.png"),
    },
    {
      "@type": "WebSite",
      "@id": `${siteOrigin}/#website`,
      name: "Referral",
      url: `${siteOrigin}/`,
      inLanguage: ["zh-CN", "en-US"],
      publisher: { "@id": `${siteOrigin}/#organization` },
    },
    {
      "@type": "SoftwareApplication",
      "@id": `${siteOrigin}/#software`,
      name: "Referral",
      applicationCategory: "BusinessApplication",
      operatingSystem: "Web",
      url: `${siteOrigin}/`,
      description:
        "A bilingual referral rewards platform for sharing invitation links, tracking successful referrals, and managing Credit rewards.",
      publisher: { "@id": `${siteOrigin}/#organization` },
      offers: {
        "@type": "Offer",
        price: "0",
        priceCurrency: "USD",
        url: absoluteURL("/register"),
      },
    },
  ],
}
