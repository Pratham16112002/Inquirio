import React from 'react'

const ActivationLayout = ({
    children
}: Readonly<{
    children: React.ReactNode
}>) => {
    return (
        <div className="flex min-h-screen items-center justify-center bg-zinc-50">
            <main className="flex flex-1  flex-col items-center justify-center">
                {children}
            </main>
        </div>
    )
}

export default ActivationLayout
