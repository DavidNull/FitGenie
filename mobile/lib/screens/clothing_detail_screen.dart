import 'package:flutter/material.dart';
import '../models/clothing_item.dart';

class ClothingDetailScreen extends StatelessWidget {
  final ClothingItem item;

  const ClothingDetailScreen({
    super.key,
    required this.item,
  });

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFE9ECF1),
      body: CustomScrollView(
        slivers: [
          SliverAppBar(
            expandedHeight: 400,
            pinned: true,
            backgroundColor: const Color(0xFF0E4A88),
            leading: IconButton(
              icon: const Icon(Icons.arrow_back, color: Colors.white),
              onPressed: () => Navigator.pop(context),
            ),
            actions: [
              IconButton(
                icon: const Icon(Icons.edit, color: Colors.white),
                onPressed: () => _editItem(context),
              ),
            ],
            flexibleSpace: FlexibleSpaceBar(
              background: _buildImage(),
            ),
          ),
          SliverToBoxAdapter(
            child: Container(
              padding: const EdgeInsets.all(24),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  _buildNameSection(),
                  const SizedBox(height: 24),
                  _buildInfoGrid(),
                  const SizedBox(height: 24),
                  _buildTagsSection(),
                  const SizedBox(height: 32),
                  _buildActionButtons(context),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildImage() {
    final color = _getCategoryColor(item.category);
    
    return Container(
      color: color.withOpacity(0.1),
      child: item.imageUrl != null && item.imageUrl!.isNotEmpty
          ? Image.network(
              item.imageUrl!,
              fit: BoxFit.cover,
              width: double.infinity,
              height: double.infinity,
              errorBuilder: (context, error, stackTrace) => _buildPlaceholder(color),
            )
          : _buildPlaceholder(color),
    );
  }

  Widget _buildPlaceholder(Color color) {
    return Center(
      child: Icon(
        _getCategoryIcon(item.category),
        size: 120,
        color: color.withOpacity(0.5),
      ),
    );
  }

  Widget _buildNameSection() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          item.name,
          style: const TextStyle(
            fontSize: 28,
            fontWeight: FontWeight.bold,
            color: Color(0xFF0E4A88),
          ),
        ),
        const SizedBox(height: 8),
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
          decoration: BoxDecoration(
            color: _getCategoryColor(item.category).withOpacity(0.2),
            borderRadius: BorderRadius.circular(20),
          ),
          child: Text(
            item.category,
            style: TextStyle(
              fontSize: 14,
              fontWeight: FontWeight.w600,
              color: _getCategoryColor(item.category),
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildInfoGrid() {
    return Column(
      children: [
        Row(
          children: [
            Expanded(
              child: _buildInfoCard(
                icon: Icons.color_lens,
                label: 'Color',
                value: item.primaryColor ?? 'No especificado',
                color: const Color(0xFFF78400),
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: _buildInfoCard(
                icon: Icons.style,
                label: 'Estilo',
                value: item.style ?? 'No especificado',
                color: const Color(0xFF1DA9B6),
              ),
            ),
          ],
        ),
        const SizedBox(height: 12),
        Row(
          children: [
            Expanded(
              child: _buildInfoCard(
                icon: Icons.event,
                label: 'Ocasión',
                value: item.occasion?.join(', ') ?? 'No especificado',
                color: const Color(0xFF0E4A88),
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: _buildInfoCard(
                icon: Icons.wb_sunny,
                label: 'Temporada',
                value: item.season?.join(', ') ?? 'No especificado',
                color: const Color(0xFFF78400),
              ),
            ),
          ],
        ),
      ],
    );
  }

  Widget _buildInfoCard({
    required IconData icon,
    required String label,
    required String value,
    required Color color,
  }) {
    return Container(
      padding: const EdgeInsets.all(16),
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
          Row(
            children: [
              Icon(icon, size: 20, color: color),
              const SizedBox(width: 8),
              Text(
                label,
                style: TextStyle(
                  fontSize: 12,
                  color: Colors.grey.shade600,
                ),
              ),
            ],
          ),
          const SizedBox(height: 8),
          Text(
            value,
            style: const TextStyle(
              fontSize: 14,
              fontWeight: FontWeight.w600,
              color: Color(0xFF0E4A88),
            ),
            maxLines: 2,
            overflow: TextOverflow.ellipsis,
          ),
        ],
      ),
    );
  }

  Widget _buildTagsSection() {
    final tags = <String>[];
    
    if (tags.isEmpty) return const SizedBox.shrink();

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'Etiquetas',
          style: TextStyle(
            fontSize: 16,
            fontWeight: FontWeight.bold,
            color: Color(0xFF0E4A88),
          ),
        ),
        const SizedBox(height: 12),
        Wrap(
          spacing: 8,
          runSpacing: 8,
          children: tags.map((tag) => _buildTag(tag)).toList(),
        ),
      ],
    );
  }

  Widget _buildTag(String label) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      decoration: BoxDecoration(
        color: const Color(0xFFF78400).withOpacity(0.2),
        borderRadius: BorderRadius.circular(20),
        border: Border.all(
          color: const Color(0xFFF78400).withOpacity(0.3),
        ),
      ),
      child: Text(
        label,
        style: const TextStyle(
          fontSize: 12,
          fontWeight: FontWeight.w600,
          color: Color(0xFFF78400),
        ),
      ),
    );
  }

  Widget _buildActionButtons(BuildContext context) {
    return Column(
      children: [
        SizedBox(
          width: double.infinity,
          child: ElevatedButton.icon(
            onPressed: () => _editItem(context),
            icon: const Icon(Icons.edit),
            label: const Text('Editar prenda'),
            style: ElevatedButton.styleFrom(
              backgroundColor: const Color(0xFF1DA9B6),
              foregroundColor: Colors.white,
              padding: const EdgeInsets.symmetric(vertical: 16),
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(12),
              ),
            ),
          ),
        ),
        const SizedBox(height: 12),
        SizedBox(
          width: double.infinity,
          child: OutlinedButton.icon(
            onPressed: () => _deleteItem(context),
            icon: const Icon(Icons.delete, color: Colors.red),
            label: const Text(
              'Eliminar prenda',
              style: TextStyle(color: Colors.red),
            ),
            style: OutlinedButton.styleFrom(
              side: const BorderSide(color: Colors.red),
              padding: const EdgeInsets.symmetric(vertical: 16),
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(12),
              ),
            ),
          ),
        ),
      ],
    );
  }

  void _editItem(BuildContext context) {
    // TODO: Navigate to edit screen
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Editar - próximamente')),
    );
  }

  void _deleteItem(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Eliminar prenda'),
        content: Text('¿Seguro que quieres eliminar "${item.name}"?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancelar'),
          ),
          TextButton(
            onPressed: () {
              Navigator.pop(context);
              Navigator.pop(context, 'deleted');
            },
            style: TextButton.styleFrom(foregroundColor: Colors.red),
            child: const Text('Eliminar'),
          ),
        ],
      ),
    );
  }

  Color _getCategoryColor(String category) {
    final lower = category.toLowerCase();
    if (lower.contains('top') || 
        lower.contains('shirt') || 
        lower.contains('camiseta') ||
        lower.contains('camisa') ||
        lower.contains('parte de arriba')) {
      return const Color(0xFF1DA9B6);
    }
    if (lower.contains('bottom') || 
        lower.contains('pant') || 
        lower.contains('parte de abajo')) {
      return const Color(0xFFF78400);
    }
    if (lower.contains('shoe') || 
        lower.contains('footwear') ||
        lower.contains('calzado')) {
      return const Color(0xFF0E4A88);
    }
    return const Color(0xFF1DA9B6);
  }

  IconData _getCategoryIcon(String category) {
    final lower = category.toLowerCase();
    if (lower.contains('top') || 
        lower.contains('shirt') || 
        lower.contains('camiseta') ||
        lower.contains('parte de arriba')) {
      return Icons.checkroom;
    }
    if (lower.contains('bottom') || 
        lower.contains('pant') || 
        lower.contains('parte de abajo')) {
      return Icons.accessibility;
    }
    if (lower.contains('shoe') || 
        lower.contains('calzado')) {
      return Icons.directions_walk;
    }
    return Icons.checkroom;
  }
}
