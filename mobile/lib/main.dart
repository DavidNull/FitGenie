import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'screens/home_screen.dart';
import 'screens/gallery_screen.dart';
import 'screens/camera_screen.dart';
import 'screens/recommendations_screen.dart';
import 'providers/app_provider.dart';
import 'services/api_service.dart';

void main() {
  WidgetsFlutterBinding.ensureInitialized();
  ApiService.detectBackend();
  runApp(
    ChangeNotifierProvider(
      create: (context) {
        final provider = AppProvider();
        provider.initialize();
        return provider;
      },
      child: const FitGenieApp(),
    ),
  );
}

class FitGenieApp extends StatelessWidget {
  const FitGenieApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'FitGenie',
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(
          seedColor: const Color(0xFF0E4A88),
          primary: const Color(0xFF0E4A88),
          secondary: const Color(0xFF1DA9B6),
          surface: const Color(0xFFE9ECF1),
        ),
        useMaterial3: true,
      ),
      home: const MainScreen(),
    );
  }
}

class MainScreen extends StatefulWidget {
  const MainScreen({super.key});

  @override
  State<MainScreen> createState() => _MainScreenState();
}

class _MainScreenState extends State<MainScreen> {
  int _selectedIndex = 0;

  late final List<Widget> _screens;

  @override
  void initState() {
    super.initState();
    _screens = [
      HomeScreen(
        onNavigateToRecommendations: () => _onItemTapped(2),
      ),
      const GalleryScreen(),
      const RecommendationsScreen(),
      const CameraScreen(),
    ];
  }

  void _onItemTapped(int index) {
    setState(() {
      _selectedIndex = index;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: _screens[_selectedIndex],
      bottomNavigationBar: Container(
        decoration: BoxDecoration(
          color: const Color(0xFF1A2A3A), // Fondo oscuro azulado para iconos sin fondo
          boxShadow: [
            BoxShadow(
              color: Colors.black.withOpacity(0.2),
              blurRadius: 15,
              offset: const Offset(0, -5),
            ),
          ],
        ),
        child: SafeArea(
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                _buildNavItem(0, 'assets/HOME.png'),
                _buildNavItem(1, 'assets/GALLERY.png'),
                _buildIconNavItem(2, Icons.auto_awesome),
                _buildNavItem(3, 'assets/camera.png'),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildNavItem(int index, String iconPath) {
    final isSelected = _selectedIndex == index;
    return GestureDetector(
      onTap: () => _onItemTapped(index),
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 200),
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: isSelected ? const Color(0xFF0E4A88) : Colors.transparent,
          shape: BoxShape.circle,
          boxShadow: isSelected
              ? [
                  BoxShadow(
                    color: const Color(0xFF0E4A88).withOpacity(0.5),
                    blurRadius: 16,
                    spreadRadius: 4,
                  ),
                ]
              : null,
        ),
        child: ClipOval(
          child: Image.asset(
            iconPath,
            width: 44,
            height: 44,
            fit: BoxFit.cover,
          ),
        ),
      ),
    );
  }

  Widget _buildIconNavItem(int index, IconData icon) {
    final isSelected = _selectedIndex == index;
    return GestureDetector(
      onTap: () => _onItemTapped(index),
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 200),
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: isSelected ? const Color(0xFF0E4A88) : Colors.transparent,
          shape: BoxShape.circle,
          boxShadow: isSelected
              ? [
                  BoxShadow(
                    color: const Color(0xFF0E4A88).withOpacity(0.5),
                    blurRadius: 16,
                    spreadRadius: 4,
                  ),
                ]
              : null,
        ),
        child: Icon(
          icon,
          size: 32,
          color: isSelected ? Colors.white : const Color(0xFF0E4A88),
        ),
      ),
    );
  }
}
