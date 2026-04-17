import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/app_provider.dart';
import '../models/clothing_item.dart';
import '../widgets/skeleton_loading.dart';
import '../utils/page_transitions.dart';
import 'clothing_detail_screen.dart';
import 'camera_screen.dart';

class GalleryScreen extends StatefulWidget {
  const GalleryScreen({super.key});

  @override
  State<GalleryScreen> createState() => _GalleryScreenState();
}

class _GalleryScreenState extends State<GalleryScreen> {
  String _selectedFilter = 'Todos';
  
  final List<String> _filters = ['Todos', 'Parte de arriba', 'Parte de abajo', 'Calzado'];

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<AppProvider>().loadClothingItems();
    });
  }
  
  List<ClothingItem> _getFilteredItems(List<ClothingItem> items) {
    if (_selectedFilter == 'Todos') return items;
    
    return items.where((item) {
      final cat = item.category.toLowerCase();
      switch (_selectedFilter) {
        case 'Parte de arriba':
          return cat.contains('top') || 
                 cat.contains('camiseta') || 
                 cat.contains('camisa') ||
                 cat.contains('chaqueta') ||
                 cat.contains('sudadera');
        case 'Parte de abajo':
          return cat.contains('bottom') || 
                 cat.contains('pantalon') || 
                 cat.contains('vaquero') ||
                 cat.contains('short') ||
                 cat.contains('falda');
        case 'Calzado':
          return cat.contains('shoe') || 
                 cat.contains('calzado') || 
                 cat.contains('zapatilla') ||
                 cat.contains('zapato');
        default:
          return true;
      }
    }).toList();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        height: MediaQuery.of(context).size.height,
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
          bottom: false, // Allow gradient to extend behind bottom nav
          child: Consumer<AppProvider>(
            builder: (context, provider, child) {
              return RefreshIndicator(
                onRefresh: () async {
                  await provider.loadClothingItems();
                },
                color: const Color(0xFF0E4A88),
                backgroundColor: Colors.white,
                child: SingleChildScrollView(
                  physics: const AlwaysScrollableScrollPhysics(),
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
                      SingleChildScrollView(
                        scrollDirection: Axis.horizontal,
                        child: Row(
                          children: _filters.map((filter) {
                            return Padding(
                              padding: const EdgeInsets.only(right: 8),
                              child: GestureDetector(
                                onTap: () => setState(() => _selectedFilter = filter),
                                child: _buildFilterChip(filter, _selectedFilter == filter),
                              ),
                            );
                          }).toList(),
                        ),
                      ),
                      const SizedBox(height: 20),
                      if (provider.isLoading && provider.clothingItems.isEmpty)
                        const SkeletonGrid(itemCount: 6)
                      else if (provider.error != null)
                        ErrorState(
                          message: provider.error!,
                          onRetry: () => provider.loadClothingItems(),
                        )
                      else if (provider.clothingItems.isEmpty)
                        Column(
                          children: [
                            EmptyState(
                              icon: Icons.checkroom,
                              title: 'Tu armario está vacío',
                              subtitle: 'Añade prendas o usa las imágenes de ejemplo',
                            ),
                            const SizedBox(height: 20),
                            ElevatedButton.icon(
                              onPressed: () => provider.importSampleImages(),
                              icon: const Icon(Icons.auto_fix_high),
                              label: const Text('Usar imágenes de ejemplo'),
                              style: ElevatedButton.styleFrom(
                                backgroundColor: const Color(0xFF1DA9B6),
                                foregroundColor: Colors.white,
                                padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
                              ),
                            ),
                            const SizedBox(height: 12),
                            OutlinedButton.icon(
                              onPressed: () {
                                PageTransitions.modalTo(
                                  context,
                                  const CameraScreen(),
                                );
                              },
                              icon: const Icon(Icons.add),
                              label: const Text('Añadir mis propias prendas'),
                              style: OutlinedButton.styleFrom(
                                foregroundColor: const Color(0xFF0E4A88),
                                side: const BorderSide(color: Color(0xFF0E4A88)),
                                padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
                              ),
                            ),
                          ],
                        )
                      else
                        Builder(
                          builder: (context) {
                            final filteredItems = _getFilteredItems(provider.clothingItems);
                            if (filteredItems.isEmpty) {
                              return EmptyState(
                                icon: Icons.filter_list,
                                title: 'No hay prendas',
                                subtitle: 'No tienes prendas en la categoría "$_selectedFilter"',
                              );
                            }
                            return LayoutBuilder(
                              builder: (context, constraints) {
                                final itemWidth = (constraints.maxWidth - 12) / 2;
                                return Wrap(
                                  spacing: 12,
                                  runSpacing: 12,
                                  children: filteredItems.map((item) {
                                    return SizedBox(
                                      width: itemWidth,
                                      child: _buildClothingItem(item),
                                    );
                                  }).toList(),
                                );
                              },
                            );
                          }
                        ),
                      ],
                    ),
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
    
    return GestureDetector(
      onTap: () => _viewClothingDetail(item),
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
          ClipRRect(
            borderRadius: const BorderRadius.vertical(
              top: Radius.circular(16),
            ),
            child: Stack(
              children: [
                Container(
                  color: color.withOpacity(0.1),
                  child: AspectRatio(
                    aspectRatio: 3 / 4,
                    child: _buildItemImage(item, color),
                  ),
                ),
                Positioned(
                  top: 8,
                  right: 8,
                  child: GestureDetector(
                    onTap: () => _deleteClothingItem(item),
                    child: Container(
                      padding: const EdgeInsets.all(6),
                      decoration: BoxDecoration(
                        color: Colors.red.withOpacity(0.9),
                        shape: BoxShape.circle,
                      ),
                      child: const Icon(
                        Icons.delete,
                        color: Colors.white,
                        size: 18,
                      ),
                    ),
                  ),
                ),
              ],
            ),
          ),
          Padding(
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
        ],
      ),
    ),
    );
  }

  void _viewClothingDetail(ClothingItem item) {
    PageTransitions.slideTo(
      context,
      ClothingDetailScreen(item: item),
    );
  }

  Future<void> _deleteClothingItem(ClothingItem item) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Eliminar prenda'),
        content: Text('¿Seguro que quieres eliminar "${item.name}"?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: const Text('Cancelar'),
          ),
          TextButton(
            onPressed: () => Navigator.of(context).pop(true),
            style: TextButton.styleFrom(foregroundColor: Colors.red),
            child: const Text('Eliminar'),
          ),
        ],
      ),
    );

    if (confirmed == true && mounted) {
      final provider = Provider.of<AppProvider>(context, listen: false);
      try {
        await provider.deleteClothingItem(item.id);
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(
              content: Text('Prenda eliminada'),
              backgroundColor: Color(0xFF1DA9B6),
            ),
          );
        }
      } catch (e) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text('Error: $e'),
              backgroundColor: Colors.red,
            ),
          );
        }
      }
    }
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
    if (item.imageUrl != null && item.imageUrl!.isNotEmpty) {
      if (item.imageUrl!.startsWith('assets/')) {
        return Image.asset(
          item.imageUrl!,
          fit: BoxFit.cover,
          width: double.infinity,
          errorBuilder: (context, error, stackTrace) {
            // Silently fail to placeholder - don't crash
            return Center(child: Icon(Icons.checkroom, size: 48, color: color));
          },
        );
      } else {
        return Image.network(
          item.imageUrl!,
          fit: BoxFit.cover,
          width: double.infinity,
          errorBuilder: (context, error, stackTrace) =>
            Center(child: Icon(Icons.checkroom, size: 48, color: color)),
        );
      }
    }
    return Center(child: Icon(Icons.checkroom, size: 48, color: color));
  }
}
