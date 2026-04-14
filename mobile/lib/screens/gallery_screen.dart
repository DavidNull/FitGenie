import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/app_provider.dart';
import '../models/clothing_item.dart';

class GalleryScreen extends StatefulWidget {
  const GalleryScreen({super.key});

  @override
  State<GalleryScreen> createState() => _GalleryScreenState();
}

class _GalleryScreenState extends State<GalleryScreen> {
  @override
  void initState() {
    super.initState();
    // Load data when screen opens
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<AppProvider>().loadClothingItems();
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        decoration: const BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topCenter,
            end: Alignment.bottomCenter,
            colors: [
              Color(0xFFF0F7F8), // Turquesa muy suave arriba
              Color(0xFFF0F7F8), // Más espacio del mismo color
              Color(0xFF0E4A88), // Azul oscuro abajo
            ],
            stops: [0.0, 0.4, 1.0], // 40% de la pantalla es color claro
          ),
        ),
        child: SafeArea(
          child: Consumer<AppProvider>(
            builder: (context, provider, child) {
              if (provider.isLoading) {
                return const Center(child: CircularProgressIndicator());
              }
              
              if (provider.error != null) {
                return Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Text('Error: ${provider.error}'),
                      ElevatedButton(
                        onPressed: () => provider.loadClothingItems(),
                        child: const Text('Reintentar'),
                      ),
                    ],
                  ),
                );
              }

              return SingleChildScrollView(
                child: Padding(
                  padding: const EdgeInsets.all(20),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Image.asset(
                            'assets/LOGO.png',
                            height: 60,
                          ),
                        ],
                      ),
                      const SizedBox(height: 20),
                      const Text(
                        'Mi armario',
                        style: TextStyle(
                          fontSize: 28,
                          fontWeight: FontWeight.bold,
                          color: Color(0xFF0E4A88),
                        ),
                      ),
                      const SizedBox(height: 8),
                      Text(
                        '${provider.clothingItems.length} prendas guardadas',
                        style: const TextStyle(
                          fontSize: 16,
                          color: Color(0xFF1DA9B6),
                        ),
                      ),
                      const SizedBox(height: 20),
                      Row(
                        children: [
                          _buildFilterChip('Todos', true),
                          const SizedBox(width: 8),
                          _buildFilterChip('Camisetas', false),
                          const SizedBox(width: 8),
                          _buildFilterChip('Pantalones', false),
                        ],
                      ),
                      const SizedBox(height: 20),
                      if (provider.clothingItems.isEmpty)
                        const Center(
                          child: Padding(
                            padding: EdgeInsets.all(40),
                            child: Text(
                              'No tienes prendas guardadas.\nAñade prendas desde la cámara.',
                              textAlign: TextAlign.center,
                              style: TextStyle(
                                fontSize: 16,
                                color: Colors.grey,
                              ),
                            ),
                          ),
                        )
                      else
                        LayoutBuilder(
                          builder: (context, constraints) {
                            final itemWidth = (constraints.maxWidth - 12) / 2;
                            return Wrap(
                              spacing: 12,
                              runSpacing: 12,
                              children: provider.clothingItems.map((item) {
                                return SizedBox(
                                  width: itemWidth,
                                  child: _buildClothingItem(item),
                                );
                              }).toList(),
                            );
                          },
                        ),
                    ],
                  ),
                ),
              );
            },
          ),
        ),
      ),
    );
  }

  Widget _buildFilterChip(String label, bool isSelected) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      decoration: BoxDecoration(
        color: isSelected ? const Color(0xFF0E4A88) : Colors.white,
        borderRadius: BorderRadius.circular(20),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 4,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: Text(
        label,
        style: TextStyle(
          color: isSelected ? Colors.white : const Color(0xFF0E4A88),
          fontWeight: isSelected ? FontWeight.bold : FontWeight.normal,
        ),
      ),
    );
  }

  Widget _buildClothingItem(ClothingItem item) {
    final color = _getCategoryColor(item.category);
    
    return SizedBox(
      height: 200,
      child: Container(
        decoration: BoxDecoration(
          color: Colors.white,
          borderRadius: BorderRadius.circular(16),
          boxShadow: [
            BoxShadow(
              color: Colors.black.withOpacity(0.05),
              blurRadius: 10,
              offset: const Offset(0, 4),
            ),
          ],
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Expanded(
              flex: 3,
              child: Container(
                decoration: BoxDecoration(
                  color: color.withOpacity(0.1),
                  borderRadius: const BorderRadius.vertical(
                    top: Radius.circular(16),
                  ),
                ),
                child: Center(
                  child: _buildItemImage(item, color),
                ),
              ),
            ),
            Expanded(
              flex: 2,
              child: Padding(
                padding: const EdgeInsets.all(12),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      item.name,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                      style: const TextStyle(
                        fontWeight: FontWeight.bold,
                        color: Color(0xFF0E4A88),
                      ),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      item.category,
                      style: const TextStyle(
                        fontSize: 12,
                        color: Colors.grey,
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Color _getCategoryColor(String category) {
    switch (category.toLowerCase()) {
      case 'camisetas':
      case 'camisas':
        return const Color(0xFF0E4A88); // Azul
      case 'pantalones':
        return const Color(0xFF1DA9B6); // Turquesa
      case 'calzado':
        return const Color(0xFFF78400); // Naranja
      default:
        return const Color(0xFF0E4A88);
    }
  }

  Widget _buildItemImage(ClothingItem item, Color color) {
    // If image URL is available, try to load it
    if (item.imageUrl != null && item.imageUrl!.isNotEmpty) {
      if (item.imageUrl!.startsWith('assets/')) {
        // Local asset image
        return Image.asset(
          item.imageUrl!,
          fit: BoxFit.cover,
          errorBuilder: (context, error, stackTrace) =>
            Icon(Icons.checkroom, size: 48, color: color),
        );
      } else {
        // Network image
        return Image.network(
          item.imageUrl!,
          fit: BoxFit.cover,
          errorBuilder: (context, error, stackTrace) =>
            Icon(Icons.checkroom, size: 48, color: color),
        );
      }
    }
    // Fallback to icon
    return Icon(Icons.checkroom, size: 48, color: color);
  }
}
