import Link from "next/link";

export default function Home() {
  return (
    <div className="min-h-screen bg-gray-50 flex flex-col font-sans">
      {/* Top Navigation Bar */}
      <header className="w-full p-6 flex justify-between items-center bg-white shadow-sm border-b border-gray-100">
        <h1 className="text-3xl font-black text-blue-600 tracking-tighter">JHINKX</h1>
        <div className="space-x-2 sm:space-x-4">
          <Link href="/signin" className="text-gray-600 hover:text-blue-600 font-medium transition-colors px-3 py-2">
            Log In
          </Link>
          <Link href="/signup" className="bg-blue-600 text-white px-5 py-2.5 rounded-lg font-bold hover:bg-blue-700 transition-colors shadow-sm">
            Sign Up
          </Link>
        </div>
      </header>

      {/* Hero Section */}
      <main className="flex-grow flex flex-col items-center justify-center text-center px-4 sm:px-6 lg:px-8 mt-12 mb-20">
        <h2 className="text-5xl md:text-7xl font-extrabold text-gray-900 tracking-tight mb-6">
          Find your perfect squad. <br className="hidden md:block" />
          <span className="text-transparent bg-clip-text bg-gradient-to-r from-blue-600 to-indigo-600">
            Instantly.
          </span>
        </h2>
        
        <p className="mt-4 max-w-2xl text-xl text-gray-500 mb-10 leading-relaxed">
          Stop scrolling through dead forums and crowded Discord servers. Tell us what you need, hit search, and drop directly into a private room with your new team.
        </p>
        
        {/* Call to Action Buttons */}
        <div className="flex flex-col sm:flex-row gap-4 w-full sm:w-auto">
          <Link href="/signup" className="w-full sm:w-auto flex items-center justify-center px-8 py-4 text-lg font-bold rounded-xl text-white bg-blue-600 hover:bg-blue-700 shadow-lg hover:shadow-xl transition-all hover:-translate-y-1">
            Find Your Team Now
          </Link>
          <Link href="/signin" className="w-full sm:w-auto flex items-center justify-center px-8 py-4 text-lg font-bold rounded-xl text-blue-700 bg-blue-50 hover:bg-blue-100 border border-blue-200 transition-colors">
            I already have an account
          </Link>
        </div>

        {/* Feature Highlights Grid */}
        <div className="mt-32 grid grid-cols-1 md:grid-cols-3 gap-8 max-w-6xl w-full text-left">
          <div className="bg-white p-8 rounded-2xl shadow-sm border border-gray-100 hover:shadow-md transition-shadow">
            <div className="w-14 h-14 bg-blue-100 text-blue-600 rounded-xl flex items-center justify-center text-3xl mb-6">⚡</div>
            <h3 className="text-xl font-bold text-gray-900 mb-3">Lightning Fast</h3>
            <p className="text-gray-500 leading-relaxed">No more waiting around. Our automated system pairs you with active users the exact moment you click search.</p>
          </div>
          
          <div className="bg-white p-8 rounded-2xl shadow-sm border border-gray-100 hover:shadow-md transition-shadow">
            <div className="w-14 h-14 bg-green-100 text-green-600 rounded-xl flex items-center justify-center text-3xl mb-6">💬</div>
            <h3 className="text-xl font-bold text-gray-900 mb-3">Instant Live Chat</h3>
            <p className="text-gray-500 leading-relaxed">The second your group is formed, you are dropped into a private, real-time chat space to strategize and connect.</p>
          </div>
          
          <div className="bg-white p-8 rounded-2xl shadow-sm border border-gray-100 hover:shadow-md transition-shadow">
            <div className="w-14 h-14 bg-purple-100 text-purple-600 rounded-xl flex items-center justify-center text-3xl mb-6">🎯</div>
            <h3 className="text-xl font-bold text-gray-900 mb-3">The Perfect Fit</h3>
            <p className="text-gray-500 leading-relaxed">Looking for a quick duo or a full 6-player raid team? You set the rules, and we build the exact squad you need.</p>
          </div>
        </div>
      </main>
    </div>
  );
}