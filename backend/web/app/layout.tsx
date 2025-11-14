import type { Metadata } from "next";
import "./globals.css"




export const metadata: Metadata = {
  title: "Account Activation",
  description: "Account activation page",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>
        {children}
      </body>
    </html>
  );
}
