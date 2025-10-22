import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "SapientIA",
  description: "Análise de conversas e dados textuais corporativos",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="pt-BR">
      <body>{children}</body>
    </html>
  );
}
